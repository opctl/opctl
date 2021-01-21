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
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

func runOp(
	ctx context.Context,
	nodeProvider nodeprovider.NodeProvider,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	cliOutput clioutput.CliOutput,
	opRef *string,
	args *[]string,
	argFile *string,
) error {
	node, err := nodeProvider.CreateNodeIfNotExists()
	if err != nil {
		return err
	}

	argMap, opHandle, err := prepareOp(ctx, *opRef, args, argFile, cliParamSatisfier, node)
	if err != nil {
		return err
	}

	return runResolvedOp(ctx, argMap, opHandle, node, cliOutput)
}

// prepareOp takes cli arguments and resolves the proper arguments and op to run
func prepareOp(
	ctx context.Context,
	opRef string,
	args *[]string,
	argFile *string,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	node node.Node,
) (map[string]*model.Value, model.DataHandle, error) {
	dataResolver := dataresolver.New(
		cliParamSatisfier,
		node,
	)

	opHandle, err := dataResolver.Resolve(
		opRef,
		nil,
	)
	if nil != err {
		return nil, nil, err
	}

	opFileReader, err := opHandle.GetContent(
		ctx,
		opfile.FileName,
	)
	if nil != err {
		return nil, nil, err
	}

	opFileBytes, err := ioutil.ReadAll(opFileReader)
	if nil != err {
		return nil, nil, err
	}

	opFile, err := opfile.Unmarshal(
		opFileBytes,
	)
	if nil != err {
		return nil, nil, err
	}

	ymlFileInputSrc, err := cliParamSatisfier.NewYMLFileInputSrc(*argFile)
	if nil != err {
		return nil, nil, fmt.Errorf("unable to load arg file at '%v'; error was: %v", *argFile, err.Error())
	}

	argsMap, err := cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			cliParamSatisfier.NewSliceInputSrc(*args, "="),
			ymlFileInputSrc,
			cliParamSatisfier.NewEnvVarInputSrc(),
			cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs),
			cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs),
		),
		opFile.Inputs,
	)
	if nil != err {
		return nil, nil, err
	}

	return argsMap, opHandle, nil
}

// runResolvedOp runs a resolved op through the CLI. It's in charge of signal handling
// and op event handling.
func runResolvedOp(
	ctx context.Context,
	argsMap map[string]*model.Value,
	opHandle model.DataHandle,
	node node.Node,
	cliOutput clioutput.CliOutput,
) error {
	startTime := time.Now().UTC()

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
	if nil != err {
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
	if nil != err {
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

			if nil != event.CallEnded {
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
