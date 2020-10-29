package core

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/parallelCaller.go . parallelCaller
type parallelCaller interface {
	// Executes a parallel call
	Call(
		parentCtx context.Context,
		callID string,
		inboundScope map[string]*model.Value,
		rootCallID string,
		opPath string,
		callSpecParallelCall []*model.CallSpec,
	)
}

func newParallelCaller(
	caller caller,
	pubSub pubsub.PubSub,
) parallelCaller {

	return _parallelCaller{
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}

}

func refToName(ref string) string {
	return strings.TrimSuffix(strings.TrimPrefix(ref, "$("), ")")
}

type _parallelCaller struct {
	caller              caller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (pc _parallelCaller) Call(
	parentCtx context.Context,
	callID string,
	inboundScope map[string]*model.Value,
	rootCallID string,
	opPath string,
	callSpecParallelCall []*model.CallSpec,
) {
	// setup cancellation
	parallelCtx, cancelParallel := context.WithCancel(parentCtx)
	defer cancelParallel()

	outputs := map[string]*model.Value{}
	var err error

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			CallEnded: &model.CallEnded{
				CallID:     callID,
				Outputs:    outputs,
				RootCallID: rootCallID,
			},
			Timestamp: time.Now().UTC(),
		}

		if nil != err {
			event.CallEnded.Error = &model.CallEndedError{
				Message: err.Error(),
			}
		}
		pc.pubSub.Publish(event)
	}()

	childCallNeededCountByName := map[string]int{}
	for _, callSpecChildCall := range callSpecParallelCall {
		// increment needed by counts for any needs
		for _, neededCallRef := range callSpecChildCall.Needs {
			childCallNeededCountByName[refToName(neededCallRef)]++
		}
	}

	startTime := time.Now().UTC()
	childCallIndexByID := map[string]int{}
	childCallIDByName := map[string]string{}
	childCallOutputsByIndex := map[int]map[string]*model.Value{}

	// perform calls in parallel w/ cancellation
	for childCallIndex, childCall := range callSpecParallelCall {

		var childCallID string
		childCallID, err = pc.uniqueStringFactory.Construct()
		if nil != err {
			// cancel all children on any error
			cancelParallel()
		}
		childCallIndexByID[childCallID] = childCallIndex

		if nil != childCall.Name {
			childCallIDByName[*childCall.Name] = childCallID
		}

		go func(childCall *model.CallSpec) {
			defer func() {
				if panicArg := recover(); panicArg != nil {
					// recover from panics; treat as errors
					err = fmt.Errorf("%v\n%v", panicArg, debug.Stack())

					// cancel all children on any error
					cancelParallel()
				}
			}()

			pc.caller.Call(
				parallelCtx,
				childCallID,
				inboundScope,
				childCall,
				opPath,
				&callID,
				rootCallID,
			)

		}(childCall)
	}

	// subscribe to events
	// @TODO: handle err channel
	eventChannel, _ := pc.pubSub.Subscribe(
		parallelCtx,
		model.EventFilter{
			Roots: []string{rootCallID},
			Since: &startTime,
		},
	)

	childErrorMessages := []string{}
	for event := range eventChannel {
		if nil != event.CallEnded {
			if childCallIndex, isChildCallEnded := childCallIndexByID[event.CallEnded.CallID]; isChildCallEnded {
				childCallOutputsByIndex[childCallIndex] = event.CallEnded.Outputs
				if nil != event.CallEnded.Error {
					// cancel all children on any error
					cancelParallel()
					childErrorMessages = append(childErrorMessages, event.CallEnded.Error.Message)
				}

				// decrement needed by counts for any needs
				for _, neededCallRef := range callSpecParallelCall[childCallIndex].Needs {
					childCallNeededCountByName[refToName(neededCallRef)]--
				}

				for neededCallName, neededCount := range childCallNeededCountByName {
					if 1 > neededCount {
						if neededCallID, ok := childCallIDByName[neededCallName]; ok {
							pc.pubSub.Publish(
								model.Event{
									OpKillRequested: &model.OpKillRequested{
										Request: model.KillOpReq{
											OpID:       neededCallID,
											RootCallID: rootCallID,
										},
									},
									Timestamp: time.Now().UTC(),
								},
							)
						}
					}
				}
			}

			if len(childCallOutputsByIndex) == len(childCallIndexByID) {
				// all calls have ended

				// construct parallel outputs
				for i := 0; i < len(callSpecParallelCall); i++ {
					callOutputs := childCallOutputsByIndex[i]
					for varName, varData := range callOutputs {
						outputs[varName] = varData
					}
				}

				// construct parallel error
				if len(childErrorMessages) != 0 {
					var formattedChildErrorMessages string
					for _, childErrorMessage := range childErrorMessages {
						formattedChildErrorMessages = fmt.Sprintf("\t-%v\n", childErrorMessage)
					}
					err = fmt.Errorf(
						"-\nError(s) during parallel call. Error(s) were:\n%v\n-",
						formattedChildErrorMessages,
					)
				}

				return
			}

		}
	}

}
