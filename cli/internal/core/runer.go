package core

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
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// RunOpts are options to run a given op through the CLI
type RunOpts struct {
	ArgFile string
	Args    []string
}

// Runer exposes the "run" command
type Runer interface {
	Run(
		ctx context.Context,
		opRef string,
		opts *RunOpts,
	) error
}

// newRuner returns an initialized "run" command
func newRuner(
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	dataResolver dataresolver.DataResolver,
	opNode node.OpNode,
) Runer {
	return _runer{
		cliOutput:         cliOutput,
		cliParamSatisfier: cliParamSatisfier,
		dataResolver:      dataResolver,
		opNode:            opNode,
	}
}

type _runer struct {
	dataResolver      dataresolver.DataResolver
	cliOutput         clioutput.CliOutput
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier
	opNode            node.OpNode
}

func (ivkr _runer) Run(
	ctx context.Context,
	opRef string,
	opts *RunOpts,
) error {
	startTime := time.Now().UTC()

	opHandle, err := ivkr.dataResolver.Resolve(
		opRef,
		nil,
	)
	if nil != err {
		return err
	}

	opFileReader, err := opHandle.GetContent(
		ctx,
		opfile.FileName,
	)
	if nil != err {
		return err
	}

	opFileBytes, err := ioutil.ReadAll(opFileReader)
	if nil != err {
		return err
	}

	opFile, err := opfile.Unmarshal(
		opFileBytes,
	)
	if nil != err {
		return err
	}

	ymlFileInputSrc, err := ivkr.cliParamSatisfier.NewYMLFileInputSrc(opts.ArgFile)
	if nil != err {
		return fmt.Errorf("unable to load arg file at '%v'; error was: %v", opts.ArgFile, err.Error())
	}

	argsMap, err := ivkr.cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			ivkr.cliParamSatisfier.NewSliceInputSrc(opts.Args, "="),
			ymlFileInputSrc,
			ivkr.cliParamSatisfier.NewEnvVarInputSrc(),
			ivkr.cliParamSatisfier.NewParamDefaultInputSrc(opFile.Inputs),
			ivkr.cliParamSatisfier.NewCliPromptInputSrc(opFile.Inputs),
		),
		opFile.Inputs,
	)
	if nil != err {
		return err
	}

	// init signal channel
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
	rootCallID, err := ivkr.opNode.StartOp(
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
	eventChannel, err := ivkr.opNode.GetEventStream(
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
				ivkr.cliOutput.Warning("Gracefully stopping... (signal Control-C again to force)")
				aSigIntWasReceivedAlready = true

				ivkr.opNode.KillOp(
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
			ivkr.cliOutput.Warning("Gracefully stopping...")

			return ivkr.opNode.KillOp(
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

			ivkr.cliOutput.Event(&event)

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
