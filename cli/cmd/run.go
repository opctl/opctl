package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/opgraph"
	"github.com/opctl/opctl/sdks/go/model"
	nodePkg "github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func newRunCmd(
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
) *cobra.Command {
	argsFlagName := "args"
	argFileFlagName := "arg-file"
	noProgressFlagName := "no-progress"
	opRefArgName := "OP_REF"

	argsFlag := []string{}
	argFileFlag := ""
	noProgressFlag := false

	runCmd := cobra.Command{
		Args: cobra.ExactArgs(1),
		Example: `# Run the op defined in the '.opspec/myOp' directory of the current working directory.
opctl run myOp

# Run the op defined in the root directory of the 'github.com/opspec-pkgs/slack.chat.post-message' git 
# repository commit tagged '1.1.0'. Pass arguments for 'apiToken', 'channelName', and 'msg' inputs.
opctl run -a apiToken="my-token" -a channelName="my-channel" -a msg="hello!" github.com/opspec-pkgs/slack.chat.post-message#1.1.0
`,
		Use: fmt.Sprintf(
			"run %s",
			opRefArgName,
		),
		Short: "Run an op",
		Long: `OP_REF can be either a 'relative/path', '/absolute/path', 'host/path/repo#tag', or 'host/path/repo#tag/path'.

If an opctl node isn't reachable, one will be started automatically. 

If auth w/ the op source fails the CLI will (re)prompt for username &
password. In non-interactive terminals, the CLI will note that it can't prompt due to being in a
non-interactive terminal and exit with a non zero exit code.

Op input args are obtained from the following sources in order:
- arg provided via -a option
- arg file
- env var
- default
- prompt

If valid input args cannot be obtained, the CLI will prompt for them. In non-interactive terminals, 
the CLI will provide details about the invalid or missing input, note that it's giving up due 
to being in a non-interactive terminal and exit with a non zero exit code.

If provided args don't meet input constraints, the CLI will (re)prompt until a valid arg is obtained.

All ops &/or images pulled will be cached.

Pulling any updates to referenced images will be attempted prior to container creation. If pulling 
an updated image fails, graceful fallback to the cached image will occur.

All containers created by opctl will be attached to an overlay network and made accessible from the
opctl node and other opctl containers by their name. Containers will be removed as they exit.
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			startTime := time.Now().UTC()
			opRef := args[0]

			np, err := local.New(nodeConfig)
			if err != nil {
				return err
			}

			node, err := np.CreateNodeIfNotExists(ctx)
			if err != nil {
				return err
			}

			dataResolver := dataresolver.New(
				cliParamSatisfier,
				node,
			)

			opHandle, err := dataResolver.Resolve(
				ctx,
				opRef,
				nil,
			)
			if err != nil {
				return err
			}

			opFile, isGoOp, err := getRunOpFile(ctx, opHandle, node)
			if err != nil {
				return err
			}

			ymlFileInputSrc, err := cliParamSatisfier.NewYMLFileInputSrc(argFileFlag)
			if err != nil {
				return fmt.Errorf("unable to load arg file at '%v': %w", argFileFlag, err)
			}

			argsSource := cliParamSatisfier.NewSliceInputSrc(argsFlag, "=")
			envSource := cliParamSatisfier.NewEnvVarInputSrc()
			defaults := cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs)
			prompt := cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs)
			inputSourcer := cliparamsatisfier.NewInputSourcer(argsSource, ymlFileInputSrc, envSource, defaults, prompt)
			if isGoOp {
				inputSourcer = goInputSourcer{provided: cliparamsatisfier.NewInputSourcer(argsSource, ymlFileInputSrc, envSource), defaults: defaults, prompt: prompt, seen: map[string]bool{}}
			}
			argsMap, err := cliParamSatisfier.Satisfy(inputSourcer, opFile.Inputs)
			if err != nil {
				return err
			}

			// init signal channels
			aSigIntWasReceivedAlready := false
			sigIntChannel := make(chan os.Signal, 1)
			defer close(sigIntChannel)
			signal.Notify(
				sigIntChannel,
				syscall.SIGINT,
			)

			sigTermChannel := make(chan os.Signal, 1)
			defer close(sigTermChannel)
			signal.Notify(
				sigTermChannel,
				syscall.SIGTERM,
			)

			sigInfoChannel := make(chan os.Signal, 1)
			defer close(sigInfoChannel)
			signal.Notify(
				sigInfoChannel,
				syscall.Signal(0x1d), // portable version of syscall.SIGINFO
			)

			// start op
			rootCallID, err := node.StartOp(
				ctx,
				model.StartOpReq{
					Args: argsMap,
					Op: model.StartOpReqOp{
						Ref: opHandle.Ref(),
					},
				},
			)
			if err != nil {
				return err
			}

			// "request animation frame" like loop to force refresh of display loading spinners
			animationFrame := make(chan bool)
			if !noProgressFlag {
				go func() {
					for {
						time.Sleep(time.Second / 10)
						animationFrame <- true
					}
				}()
			}

			var state opgraph.CallGraph
			var loadingSpinner opgraph.DotLoadingSpinner
			output := opgraph.NewOutputManager()

			defer func() {
				output.Print(state.String(loadingSpinner, time.Now(), false))
				fmt.Println()
			}()

			clearGraph := func() {
				if !noProgressFlag {
					output.Clear()
				}
			}

			displayGraph := func() {
				if !noProgressFlag {
					output.Print(state.String(loadingSpinner, time.Now(), true))
				}
			}

			// start event loop
			eventChannel, err := node.GetEventStream(
				ctx,
				&model.GetEventStreamReq{
					Filter: model.EventFilter{
						Roots: []string{rootCallID},
						Since: &startTime,
					},
				},
			)
			if err != nil {
				return fmt.Errorf("error getting event stream: %w", err)
			}

			for {
				select {
				case <-sigIntChannel:
					clearGraph()
					if !aSigIntWasReceivedAlready {
						cliOutput.Warning("Gracefully stopping... (signal Control-C again to force)")
						aSigIntWasReceivedAlready = true

						node.KillOp(
							ctx,
							model.KillOpReq{
								OpID:       rootCallID,
								RootCallID: rootCallID,
							},
						)

						// events will continue to stream in, make sure we continue to display the graph while this happens
						displayGraph()
					} else {
						return clioutput.RunError{
							ExitCode: 130,
							Message:  "Terminated by Control-C",
						}
					}

				case <-sigInfoChannel:
					clearGraph()
					// clear two more lines
					fmt.Print("\033[1A\033[K\033[1A\033[K")
					fmt.Println(state.String(opgraph.StaticLoadingSpinner{}, time.Now(), false))
					displayGraph()

				case <-sigTermChannel:
					clearGraph()
					cliOutput.Error("Gracefully stopping...")
					node.KillOp(
						ctx,
						model.KillOpReq{
							OpID:       rootCallID,
							RootCallID: rootCallID,
						},
					)
					displayGraph()

				case event, isEventChannelOpen := <-eventChannel:
					clearGraph()
					if !isEventChannelOpen {
						return errors.New("event channel closed unexpectedly")
					}

					if err := state.HandleEvent(&event); err != nil {
						cliOutput.Error(fmt.Sprintf("%v", err))
					}

					if done, err := reportRunEvent(cliOutput, &event, rootCallID); done {
						return err
					}
					displayGraph()
				case <-animationFrame:
					clearGraph()
					displayGraph()
				}
			}
		},
	}

	runCmd.Flags().StringArrayVarP(&argsFlag, argsFlagName, "a", []string{}, "Explicitly pass args to the op")
	runCmd.Flags().StringVarP(&argFileFlag, argFileFlagName, "", filepath.Join(model.DotOpspecDirName, "args.yml"), "Read in a file of args in yml format")
	runCmd.Flags().BoolVarP(&noProgressFlag, noProgressFlagName, "", !term.IsTerminal(int(os.Stdout.Fd())), "Disable live call graph for the op")

	return &runCmd
}

func getRunOpFile(ctx context.Context, opHandle model.DataHandle, apiNode nodePkg.Node) (*model.OpSpec, bool, error) {
	isGoOp, err := opfile.IsGoOp(ctx, opHandle)
	if err != nil {
		return nil, false, err
	}
	if isGoOp {
		spec, err := apiNode.DescribeOp(ctx, model.DescribeOpReq{Op: model.StartOpReqOp{Ref: opHandle.Ref()}})
		return spec, true, err
	}
	spec, err := opfile.Get(ctx, opHandle)
	return spec, false, err
}

// reportRunEvent forwards each event before handling the root's terminal outcome.
func reportRunEvent(output clioutput.CliOutput, event *model.Event, rootID string) (bool, error) {
	output.Event(event)
	if event.CallEnded == nil || event.CallEnded.Call.ID != rootID {
		return false, nil
	}
	switch event.CallEnded.Outcome {
	case model.OpOutcomeSucceeded:
		return true, nil
	case model.OpOutcomeKilled:
		return true, clioutput.RunError{ExitCode: 137}
	default:
		return true, clioutput.RunError{ExitCode: 1}
	}
}
