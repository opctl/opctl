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
		scgParallelCall []model.NamedSCG,
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

func refToVariable(ref string) string {
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
	scgParallelCall []model.NamedSCG,
) {
	outputs := map[string]*model.Value{}
	for varName, varData := range inboundScope {
		outputs[varName] = varData
	}
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

	startTime := time.Now().UTC()

	childCallIndexByID := map[string]int{}
	for scgCallIndex, scgNamedCall := range scgParallelCall {
		for scgCallName := range scgNamedCall {
			var childCallID string
			childCallID, err = pc.uniqueStringFactory.Construct()
			if nil != err {
				return
			}
			childCallIndexByID[childCallID] = scgCallIndex

			if "" != scgCallName {
				// add named calls to scope
				outputs[scgCallName] = &model.Value{
					String: &childCallID,
				}
			}
		}
	}

	neededByCountByID := map[string]int{}
	for _, scgNamedCall := range scgParallelCall {
		for _, scgUnNamedCall := range scgNamedCall {
			// increment needed by counts for any needs
			for _, need := range scgUnNamedCall.Needs {
				neededByCountByID[*outputs[refToVariable(need)].String]++
			}
		}
	}

	// setup cancellation
	ctxOfChildren, cancelChildren := context.WithCancel(ctx)
	defer cancelChildren()

	for childCallID, childCallIndex := range childCallIndexByID {
		for _, scgUnNamedCall := range scgParallelCall[childCallIndex] {
			// loop vars same address each loop; need to copy
			scgUnNamedCall := scgUnNamedCall

			go pc.caller.Call(
				ctxOfChildren,
				childCallID,
				inboundScope,
				&scgUnNamedCall,
				opPath,
				&callID,
				rootOpID,
			)
		}
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
	callIndexOutputsMap := map[int]map[string]*model.Value{}
	for event := range eventChannel {
		if nil != event.CallEnded {
			if childCallIndex, isChildCall := childCallIndexByID[event.CallEnded.CallID]; isChildCall {
				callIndexOutputsMap[childCallIndex] = event.CallEnded.Outputs
				if nil != event.CallEnded.Error {
					// cancel all children on any error
					cancelChildren()
					childErrorMessages = append(childErrorMessages, event.CallEnded.Error.Message)
				}

				for _, scgUnNamedCall := range scgParallelCall[childCallIndex] {
					// decrement needed by counts for any needs
					for _, need := range scgUnNamedCall.Needs {
						neededByCountByID[*outputs[refToVariable(need)].String]--
					}

					for neededCallID, neededByCount := range neededByCountByID {
						if 1 > neededByCount {
							pc.callKiller.Kill(neededCallID, rootOpID)
						}
					}
				}
			}

			if len(callIndexOutputsMap) == len(childCallIndexByID) {
				// all calls have ended

				// construct parallel outputs
				for i := 0; i < len(scgParallelCall); i++ {
					callOutputs := callIndexOutputsMap[i]
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
