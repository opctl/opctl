package core

import (
	"fmt"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"os/signal"
	"path"
	"syscall"
)

func (this _core) RunOp(
	args []string,
	collection string,
	name string,
) {

	// ensure node running
	nodes, err := this.nodeProvider.ListNodes()
	if nil != err {
		panic(err.Error())
	}
	if len(nodes) < 1 {
		this.nodeProvider.CreateNode()
	}

	pwd, err := this.vos.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	opPath := path.Join(pwd, collection, name)
	opView, err := this.pkg.GetOp(opPath)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	argsMap := this.cliParamSatisfier.Satisfy(args, opView.Inputs)

	// init signal channel
	intSignalsReceived := 0
	signalChannel := make(chan os.Signal, 1)
	defer close(signalChannel)

	signal.Notify(
		signalChannel,
		syscall.SIGINT, //handle SIGINTs
	)

	// start op
	rootOpId, err := this.consumeNodeApi.StartOp(
		model.StartOpReq{
			Args:     argsMap,
			OpPkgRef: opPath,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	// start event loop
	eventChannel, err := this.consumeNodeApi.GetEventStream(
		&model.GetEventStreamReq{
			Filter: &model.EventFilter{
				RootOpIds: []string{rootOpId},
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

				this.consumeNodeApi.KillOp(
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
