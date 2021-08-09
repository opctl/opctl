package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// run implements "run" command
func run(
	ctx context.Context,
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeProvider nodeprovider.NodeProvider,
	args []string,
	argFile string,
	opRef string,
) error {

	startTime := time.Now().UTC()

	node, err := nodeProvider.CreateNodeIfNotExists(ctx)
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

	opFileBytes, err := ioutil.ReadAll(opFileReader)
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
		return err
	}

	for {
		select {

		case <-sigIntChannel:
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
			} else {
				return &RunError{
					ExitCode: 130,
					message:  "Terminated by Control-C",
				}
			}

		case <-sigTermChannel:
			cliOutput.Warning("Gracefully stopping...")

			return node.KillOp(
				ctx,
				model.KillOpReq{
					OpID:       rootCallID,
					RootCallID: rootCallID,
				},
			)
		case event, isEventChannelOpen := <-eventChannel:
			if !isEventChannelOpen {
				return errors.New("Event channel closed unexpectedly")
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
		}
	}
}
