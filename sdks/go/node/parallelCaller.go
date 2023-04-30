package node

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
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
	) (
		map[string]*model.Value,
		error,
	)
}

func newParallelCaller(
	caller caller,
	pubSub pubsub.PubSub,
) parallelCaller {

	return _parallelCaller{
		caller: caller,
		pubSub: pubSub,
	}

}

type _parallelCaller struct {
	caller caller
	pubSub pubsub.PubSub
}

func (pc _parallelCaller) Call(
	parentCtx context.Context,
	callID string,
	inboundScope map[string]*model.Value,
	rootCallID string,
	opPath string,
	callSpecParallelCall []*model.CallSpec,
) (
	map[string]*model.Value,
	error,
) {
	// setup cancellation
	parallelCtx, cancelParallel := context.WithCancel(parentCtx)
	defer cancelParallel()

	childCallNeededCountByName := map[string]int{}
	for _, callSpecChildCall := range callSpecParallelCall {
		// increment needed by counts for any needs
		for _, neededCallRef := range callSpecChildCall.Needs {
			childCallNeededCountByName[opspec.RefToName(neededCallRef)]++
		}
	}

	startTime := time.Now().UTC()
	childCallIndexByID := map[string]int{}
	childCallIDByName := map[string]string{}
	childCallOutputsByIndex := map[int]map[string]*model.Value{}

	// perform calls in parallel w/ cancellation
	for childCallIndex, childCall := range callSpecParallelCall {

		childCallID, err := uniquestring.Construct()
		if err != nil {
			// end run immediately on any error
			return nil, err
		}

		childCallIndexByID[childCallID] = childCallIndex

		if childCall.Name != nil {
			childCallIDByName[*childCall.Name] = childCallID
		}

		go func(childCall *model.CallSpec) {
			defer func() {
				if panic := recover(); panic != nil {
					// recover from panics; treat as errors
					fmt.Printf("recovered from panic: %s\n%s\n", panic, string(debug.Stack()))

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
		// don't cancel w/ children; we need to read err msgs
		parentCtx,
		model.EventFilter{
			Roots: []string{rootCallID},
			Since: &startTime,
		},
	)

	var isChildErred = false
	outputs := map[string]*model.Value{}

eventLoop:
	for event := range eventChannel {
		if event.CallEnded != nil {
			if childCallIndex, isChildCallEnded := childCallIndexByID[event.CallEnded.Call.ID]; isChildCallEnded {
				childCallOutputsByIndex[childCallIndex] = event.CallEnded.Outputs
				if event.CallEnded.Error != nil {
					isChildErred = true

					// cancel all children on any error
					cancelParallel()
				}

				// decrement needed by counts for any needs
				for _, neededCallRef := range callSpecParallelCall[childCallIndex].Needs {
					childCallNeededCountByName[opspec.RefToName(neededCallRef)]--
				}

				for neededCallName, neededCount := range childCallNeededCountByName {
					if 1 > neededCount {
						if neededCallID, ok := childCallIDByName[neededCallName]; ok {
							pc.pubSub.Publish(
								model.Event{
									CallKillRequested: &model.CallKillRequested{
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

				if isChildErred {
					return nil, errors.New("child call failed")
				}

				break eventLoop
			}

		}
	}

	return outputs, nil
}
