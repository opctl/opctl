package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/opgraph"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// run implements "run" command
func run(
	ctx context.Context,
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeConfig local.NodeConfig,
	args []string,
	argFile string,
	opRef string,
	noProgress bool,
) error {
	startTime := time.Now().UTC()

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

	ymlFileInputSrc, err := cliParamSatisfier.NewYMLFileInputSrc(argFile)
	if err != nil {
		return fmt.Errorf("unable to load arg file at '%v': %w", argFile, err)
	}

	argsMap, err := cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			cliParamSatisfier.NewSliceInputSrc(args, "="),
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
	if !noProgress {
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
		if !noProgress {
			output.Clear()
		}
	}

	displayGraph := func() {
		if !noProgress {
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
				return &RunError{
					ExitCode: 130,
					message:  "Terminated by Control-C",
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

			cliOutput.Event(&event)
			if event.CallEnded != nil {
				if event.CallEnded.Call.ID == rootCallID {
					switch event.CallEnded.Outcome {
					case model.OpOutcomeSucceeded:
						return nil
					case model.OpOutcomeKilled:
						return &RunError{ExitCode: 137}
					default:
						return &RunError{ExitCode: 1}
					}
				}
			}
			displayGraph()
		case <-animationFrame:
			clearGraph()
			displayGraph()
		}
	}
}
