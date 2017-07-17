package core

import (
	"context"
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/pkg"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (this _core) Run(
	ctx context.Context,
	pkgRef string,
	opts *RunOpts,
) {

	// ensure node running
	nodes, err := this.nodeProvider.ListNodes()
	if nil != err {
		panic(err.Error())
	}
	if len(nodes) < 1 {
		this.nodeProvider.CreateNode()
	}

	cwd, err := this.os.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	startTime := time.Now().UTC()

	pkgHandle, err := this.pkg.Resolve(pkgRef, &pkg.ResolveOpts{BasePath: cwd})
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf("Unable to resolve package '%v' from '%v'; error was: %v", pkgRef, cwd, err),
			Code:    1})
		return // support fake exiter
	}

	pkgManifest, err := this.pkg.GetManifest(pkgHandle)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	ymlFileInputSrc, err := this.cliParamSatisfier.NewYMLFileInputSrc(opts.ArgFile)
	if nil != err {
		err = fmt.Errorf("Unable to load arg file at '%v'; error was: %v", opts.ArgFile, err.Error())
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	argsMap := this.cliParamSatisfier.Satisfy(
		cliparamsatisfier.NewInputSourcer(
			this.cliParamSatisfier.NewSliceInputSrc(opts.Args, "="),
			ymlFileInputSrc,
			this.cliParamSatisfier.NewEnvVarInputSrc(),
			this.cliParamSatisfier.NewParamDefaultInputSrc(pkgManifest.Inputs),
			this.cliParamSatisfier.NewCliPromptInputSrc(pkgManifest.Inputs),
		),
		pkgManifest.Inputs,
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
	rootOpId, err := this.opspecNodeAPIClient.StartOp(
		ctx,
		model.StartOpReq{
			Args: argsMap,
			Pkg: &model.DCGOpCallPkg{
				Ref: pkgHandle.Ref(),
			},
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	// start event loop
	eventChannel, err := this.opspecNodeAPIClient.GetEventStream(
		&model.GetEventStreamReq{
			Filter: &model.EventFilter{
				RootOpIds: []string{rootOpId},
				Since:     &startTime,
			},
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {
		select {

		case <-signalChannel:
			if intSignalsReceived == 0 {

				intSignalsReceived++
				fmt.Println(this.cliColorer.Error("Gracefully stopping... (signal Control-C again to force)"))

				this.opspecNodeAPIClient.KillOp(
					ctx,
					model.KillOpReq{
						OpId: rootOpId,
					},
				)
			} else {
				this.cliExiter.Exit(cliexiter.ExitReq{Message: "Terminated by Control-C", Code: 130})
				return // support fake exiter
			}

		case event, isEventChannelOpen := <-eventChannel:
			if !isEventChannelOpen {
				this.cliExiter.Exit(cliexiter.ExitReq{Message: "Event channel closed unexpectedly", Code: 1})
				return // support fake exiter
			}

			this.cliOutput.Event(&event)

			if nil != event.OpEnded {
				if event.OpEnded.OpId == rootOpId {
					switch event.OpEnded.Outcome {
					case model.OpOutcomeSucceeded:
						this.cliExiter.Exit(cliexiter.ExitReq{Code: 0})
					case model.OpOutcomeKilled:
						this.cliExiter.Exit(cliexiter.ExitReq{Code: 137})
					default:
						// treat model.OpOutcomeFailed & unexpected values as errors.
						this.cliExiter.Exit(cliexiter.ExitReq{Code: 1})
					}
					return // support fake exiter
				}
			}
		}

	}

}
