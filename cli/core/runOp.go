package core

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
)

func (this _core) RunOp(
	args []string,
	collection string,
	name string,
) {
	pwd, err := this.vos.Getwd()
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	opPath := path.Join(filepath.Join(pwd, collection, name))
	opView, err := this.bundle.GetOp(opPath)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	argsMap := this.paramSatisfier.Satisfy(args, opView.Inputs)

	// init signal channel
	intSignalsReceived := 0
	signalChannel := make(chan os.Signal, 1)
	defer close(signalChannel)

	signal.Notify(
		signalChannel,
		syscall.SIGINT, //handle SIGINTs
	)

	// start op
	opGraphId, err := this.engineClient.StartOp(
		model.StartOpReq{
			Args:  argsMap,
			OpRef: opPath,
		},
	)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	// start event loop
	eventChannel, err := this.engineClient.GetEventStream(
		&model.GetEventStreamReq{
			Filter: &model.EventFilter{
				OpGraphIds: []string{opGraphId},
			},
		},
	)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for {
		select {

		case <-signalChannel:
			if intSignalsReceived == 0 {

				intSignalsReceived++
				fmt.Println(this.colorer.Error("Gracefully stopping... (signal Control-C again to force)"))

				this.engineClient.KillOp(
					model.KillOpReq{
						OpGraphId: opGraphId,
					},
				)
			} else {
				this.exiter.Exit(ExitReq{Message: "Terminated by Control-C", Code: 130})
				return // support fake exiter
			}

		case event, isEventChannelOpen := <-eventChannel:
			if !isEventChannelOpen {
				this.exiter.Exit(ExitReq{Message: "Event channel closed unexpectedly", Code: 1})
				return // support fake exiter
			}

			this.output.Event(&event)

			if nil != event.OpEnded {
				if event.OpEnded.OpId == opGraphId {
					switch event.OpEnded.Outcome {
					case model.OpOutcomeSucceeded:
						this.exiter.Exit(ExitReq{Code: 0})
					case model.OpOutcomeKilled:
						this.exiter.Exit(ExitReq{Code: 137})
					default:
						// treat model.OpOutcomeFailed & unexpected values as errors.
						this.exiter.Exit(ExitReq{Code: 1})
					}
					return // support fake exiter
				}
			}
		}

	}

}
