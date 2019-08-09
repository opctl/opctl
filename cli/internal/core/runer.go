package core

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opctl/opctl/cli/internal/apireachabilityensurer"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	cliModel "github.com/opctl/opctl/cli/internal/model"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/api/client"
	dotyml "github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// Runer exposes the "run" command
type Runer interface {
	Run(
		ctx context.Context,
		opRef string,
		opts *cliModel.RunOpts,
	)
}

// newRuner returns an initialized "run" command
func newRuner(
	apiClient client.Client,
	apiReachabilityEnsurer apireachabilityensurer.APIReachabilityEnsurer,
	cliColorer clicolorer.CliColorer,
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	dataResolver dataresolver.DataResolver,
) Runer {
	return _runer{
		apiClient:              apiClient,
		apiReachabilityEnsurer: apiReachabilityEnsurer,
		cliColorer:             cliColorer,
		cliExiter:              cliExiter,
		cliOutput:              cliOutput,
		cliParamSatisfier:      cliParamSatisfier,
		dataResolver:           dataResolver,
		opDotYmlGetter:         dotyml.NewGetter(),
	}
}

type _runer struct {
	apiClient              client.Client
	dataResolver           dataresolver.DataResolver
	apiReachabilityEnsurer apireachabilityensurer.APIReachabilityEnsurer
	cliColorer             clicolorer.CliColorer
	cliExiter              cliexiter.CliExiter
	cliOutput              clioutput.CliOutput
	cliParamSatisfier      cliparamsatisfier.CLIParamSatisfier
	opDotYmlGetter         dotyml.Getter
}

func (ivkr _runer) Run(
	ctx context.Context,
	opRef string,
	opts *cliModel.RunOpts,
) {
	startTime := time.Now().UTC()

	opHandle := ivkr.dataResolver.Resolve(
		opRef,
		nil,
	)

	opDotYml, err := ivkr.opDotYmlGetter.Get(
		ctx,
		opHandle,
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	ymlFileInputSrc, err := ivkr.cliParamSatisfier.NewYMLFileInputSrc(opts.ArgFile)
	if nil != err {
		err = fmt.Errorf("unable to load arg file at '%v'; error was: %v", opts.ArgFile, err.Error())
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	argsMap := ivkr.cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			ivkr.cliParamSatisfier.NewSliceInputSrc(opts.Args, "="),
			ymlFileInputSrc,
			ivkr.cliParamSatisfier.NewEnvVarInputSrc(),
			ivkr.cliParamSatisfier.NewParamDefaultInputSrc(opDotYml.Inputs),
			ivkr.cliParamSatisfier.NewCliPromptInputSrc(opDotYml.Inputs),
		),
		opDotYml.Inputs,
	)

	// init signal channel
	intSignalsReceived := 0
	signalChannel := make(chan os.Signal, 1)
	defer close(signalChannel)

	signal.Notify(
		signalChannel,
		syscall.SIGINT, //handle SIGINTs
	)

	// start op
	rootOpID, err := ivkr.apiClient.StartOp(
		ctx,
		model.StartOpReq{
			Args: argsMap,
			Op: model.StartOpReqOp{
				Ref: opHandle.Ref(),
			},
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	// start event loop
	eventChannel, err := ivkr.apiClient.GetEventStream(
		ctx,
		&model.GetEventStreamReq{
			Filter: model.EventFilter{
				Roots: []string{rootOpID},
				Since: &startTime,
			},
		},
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {
		select {

		case <-signalChannel:
			if intSignalsReceived == 0 {

				intSignalsReceived++
				fmt.Println(ivkr.cliColorer.Error("Gracefully stopping... (signal Control-C again to force)"))

				ivkr.apiClient.KillOp(
					ctx,
					model.KillOpReq{
						OpID: rootOpID,
					},
				)
			} else {
				ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: "Terminated by Control-C", Code: 130})
				return // support fake exiter
			}

		case event, isEventChannelOpen := <-eventChannel:
			if !isEventChannelOpen {
				ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1})
				return // support fake exiter
			}

			ivkr.cliOutput.Event(&event)

			if nil != event.OpEnded {
				if event.OpEnded.OpID == rootOpID {
					switch event.OpEnded.Outcome {
					case model.OpOutcomeSucceeded:
						ivkr.cliExiter.Exit(cliexiter.ExitReq{Code: 0})
					case model.OpOutcomeKilled:
						ivkr.cliExiter.Exit(cliexiter.ExitReq{Code: 137})
					default:
						// treat model.OpOutcomeFailed & unexpected values as errors.
						ivkr.cliExiter.Exit(cliexiter.ExitReq{Code: 1})
					}
					return // support fake exiter
				}
			}
		}

	}

}
