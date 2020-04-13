package core

import (
	"context"
	"fmt"
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
		ctx context.Context,
		callID string,
		inboundScope map[string]*model.Value,
		rootOpID string,
		opPath string,
		scgParallelCall []*model.SCG,
	)
}

func newParallelCaller(
	callKiller callKiller,
	caller caller,
	pubSub pubsub.PubSub,
) parallelCaller {

	return _parallelCaller{
		callKiller:          callKiller,
		caller:              caller,
		pubSub:              pubSub,
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
	}

}

func refToName(ref string) string {
	return strings.TrimSuffix(strings.TrimPrefix(ref, "$("), ")")
}

type _parallelCaller struct {
	callKiller          callKiller
	caller              caller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (pc _parallelCaller) Call(
	ctx context.Context,
	callID string,
	inboundScope map[string]*model.Value,
	rootOpID string,
	opPath string,
	scgParallelCall []*model.SCG,
) {
	// setup cancellation
	ctxOfChildren, cancelChildren := context.WithCancel(ctx)
	defer cancelChildren()

	outputs := map[string]*model.Value{}
	var err error

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		event := model.Event{
			ParallelCallEnded: &model.ParallelCallEndedEvent{
				CallID:   callID,
				Outputs:  outputs,
				RootOpID: rootOpID,
			},
			Timestamp: time.Now().UTC(),
		}

		if nil != err {
			event.ParallelCallEnded.Error = &model.CallEndedEventError{
				Message: err.Error(),
			}
		}
		pc.pubSub.Publish(
			event,
		)
	}()

	childCallNeededCountByName := map[string]int{}
	for _, scgChildCall := range scgParallelCall {
		// increment needed by counts for any needs
		for _, neededCallRef := range scgChildCall.Needs {
			childCallNeededCountByName[refToName(neededCallRef)]++
		}
	}

	startTime := time.Now().UTC()
	childCallIndexByID := map[string]int{}
	childCallIDByName := map[string]string{}
	childCallOutputsByIndex := map[int]map[string]*model.Value{}

	// perform calls in parallel w/ cancellation
	for childCallIndex, childCall := range scgParallelCall {

		var childCallID string
		childCallID, err = pc.uniqueStringFactory.Construct()
		if nil != err {
			// cancel all children on any error
			cancelChildren()
		}
		childCallIndexByID[childCallID] = childCallIndex

		if nil != childCall.Name {
			childCallIDByName[*childCall.Name] = childCallID
		}

		go pc.caller.Call(
			ctxOfChildren,
			childCallID,
			inboundScope,
			childCall,
			opPath,
			&callID,
			rootOpID,
		)
	}

	// subscribe to events
	// @TODO: handle err channel
	eventChannel, _ := pc.pubSub.Subscribe(
		ctx,
		model.EventFilter{
			Roots: []string{rootOpID},
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
					cancelChildren()
					childErrorMessages = append(childErrorMessages, event.CallEnded.Error.Message)
				}

				// decrement needed by counts for any needs
				for _, neededCallRef := range scgParallelCall[childCallIndex].Needs {
					childCallNeededCountByName[refToName(neededCallRef)]--
				}

				for neededCallName, neededCount := range childCallNeededCountByName {
					if 1 > neededCount {
						if neededCallID, ok := childCallIDByName[neededCallName]; ok {
							pc.callKiller.Kill(neededCallID, rootOpID)
						}
					}
				}
			}

			if len(childCallOutputsByIndex) == len(childCallIndexByID) {
				// all calls have ended

				// construct parallel outputs
				for i := 0; i < len(scgParallelCall); i++ {
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
