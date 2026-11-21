package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/opctl/opctl/cli/internal/attach"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/opgraph"
	"github.com/opctl/opctl/sdks/go/model"
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
			interactive := attach.IsInteractive(os.Stdin, os.Stdout)
			if interactive {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				defer cancel()
			}
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

			opFileReader, err := opHandle.GetContent(
				ctx,
				opfile.FileName,
			)
			if err != nil {
				return err
			}

			opFileBytes, err := io.ReadAll(opFileReader)
			if err != nil {
				return err
			}

			opFile, err := opfile.Unmarshal(
				opFileBytes,
			)
			if err != nil {
				return err
			}

			ymlFileInputSrc, err := cliParamSatisfier.NewYMLFileInputSrc(argFileFlag)
			if err != nil {
				return fmt.Errorf("unable to load arg file at '%v': %w", argFileFlag, err)
			}

			argsMap, err := cliParamSatisfier.Satisfy(
				cliparamsatisfier.NewInputSourcer(
					cliParamSatisfier.NewSliceInputSrc(argsFlag, "="),
					ymlFileInputSrc,
					cliParamSatisfier.NewEnvVarInputSrc(),
					cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs),
					cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs),
				),
				opFile.Inputs,
			)
			if err != nil {
				return err
			}

			// init signal channels
			aSigIntWasReceivedAlready := false
			sigIntChannel := make(chan os.Signal, 1)
			defer signal.Stop(sigIntChannel)
			signal.Notify(
				sigIntChannel,
				syscall.SIGINT,
			)

			sigTermChannel := make(chan os.Signal, 1)
			defer signal.Stop(sigTermChannel)
			signal.Notify(
				sigTermChannel,
				syscall.SIGTERM,
			)

			var sigInfoChannel, sigWinchChannel chan os.Signal
			if interactive {
				sigWinchChannel = make(chan os.Signal, 1)
				defer signal.Stop(sigWinchChannel)
				signal.Notify(sigWinchChannel, syscall.SIGWINCH)
			} else {
				sigInfoChannel = make(chan os.Signal, 1)
				defer signal.Stop(sigInfoChannel)
				signal.Notify(sigInfoChannel, syscall.Signal(0x1d)) // portable version of syscall.SIGINFO
			}

			var termSize *model.TermSize
			if interactive {
				if cols, rows, err := term.GetSize(int(os.Stdout.Fd())); err == nil && rows > 0 && cols > 0 {
					termSize = &model.TermSize{Rows: uint(rows), Cols: uint(cols)}
				}
			}

			// start op
			rootCallID, err := node.StartOp(
				ctx,
				model.StartOpReq{
					Args:        argsMap,
					Interactive: interactive,
					TermSize:    termSize,
					Op: model.StartOpReqOp{
						Ref: opHandle.Ref(),
					},
				},
			)
			if err != nil {
				return err
			}

			var animationFrame <-chan time.Time
			if !noProgressFlag {
				ticker := time.NewTicker(time.Second / 10)
				defer ticker.Stop()
				animationFrame = ticker.C
			}

			var state opgraph.CallGraph
			var loadingSpinner opgraph.DotLoadingSpinner
			output := opgraph.NewOutputManager()

			defer func() {
				output.Print(state.String(loadingSpinner, time.Now(), false))
				fmt.Println()
			}()

			runOutput := cliOutput
			var terminal *attach.Terminal
			containerFocused := func() bool { return terminal != nil && terminal.Focused() }
			graphVisible := false
			clearGraph := func() {
				if graphVisible {
					output.Clear()
					graphVisible = false
				}
			}
			displayGraph := func() {
				if showLiveGraph(noProgressFlag, containerFocused()) {
					graphVisible = output.Print(state.String(loadingSpinner, time.Now(), true)) == nil
				}
			}

			var inputTick <-chan time.Time
			if interactive {
				terminal = attach.NewTerminal(ctx, node, os.Stdin, os.Stdout, os.Stderr)
				defer terminal.Close()
				runOutput = clioutput.NewInteractive(cliColorer, terminal.MessageWriter(os.Stderr), terminal.MessageWriter(os.Stdout), terminal.WriteOutput)
				ticker := time.NewTicker(15 * time.Millisecond)
				defer ticker.Stop()
				inputTick = ticker.C
			}
			reportInputError := func(err error) {
				if err != nil {
					runOutput.Warning(err.Error())
				}
			}

			// start event loop
			eventChannel, err := node.GetEventStream(
				ctx,
				&model.GetEventStreamReq{
					Filter: model.EventFilter{
						Roots:   []string{rootCallID},
						Since:   &startTime,
						Ordered: interactive,
					},
				},
			)
			if err != nil {
				return fmt.Errorf("error getting event stream: %w", err)
			}

			for {
				select {
				case <-inputTick:
					wasFocused := containerFocused()
					reportInputError(terminal.PollInput(clearGraph))
					if wasFocused && !containerFocused() {
						displayGraph()
					}
				case <-sigWinchChannel:
					reportInputError(terminal.Resize())
				case <-sigIntChannel:
					if terminal != nil {
						reportInputError(terminal.Close())
					}
					clearGraph()
					if !aSigIntWasReceivedAlready {
						runOutput.Warning("Gracefully stopping... (signal Control-C again to force)")
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
					if terminal != nil {
						reportInputError(terminal.Close())
					}
					clearGraph()
					runOutput.Error("Gracefully stopping...")
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
						runOutput.Error(fmt.Sprintf("%v", err))
					}

					runOutput.Event(&event)
					if terminal != nil && event.ContainerStarted != nil {
						reportInputError(terminal.Started(event.ContainerStarted))
					}
					if event.CallEnded != nil {
						if container := event.CallEnded.Call.Container; terminal != nil && container != nil {
							reportInputError(terminal.Ended(container.ContainerID))
						}
						if event.CallEnded.Call.ID == rootCallID {
							switch event.CallEnded.Outcome {
							case model.OpOutcomeSucceeded:
								return nil
							case model.OpOutcomeKilled:
								return clioutput.RunError{ExitCode: 137}
							default:
								return clioutput.RunError{ExitCode: 1}
							}
						}
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

func showLiveGraph(noProgress, containerFocused bool) bool {
	return !noProgress && !containerFocused
}
