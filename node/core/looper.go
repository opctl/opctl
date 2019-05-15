package core

//go:generate counterfeiter -o ./fakeLooper.go --fake-name fakeLooper ./ looper

import (
	"context"
	"sort"
	"time"

	"github.com/opctl/sdk-golang/opspec/interpreter/loopable"

	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

type looper interface {
	// Loop loops a call
	Loop(
		ctx context.Context,
		id string,
		inboundScope map[string]*model.Value,
		scg *model.SCG,
		opHandle model.DataHandle,
		parentCallID *string,
		rootOpID string,
	) error
}

func newLooper(
	caller caller,
	pubSub pubsub.PubSub,
) looper {
	return _looper{
		caller:              caller,
		loopableInterpreter: loopable.NewInterpreter(),
		pubSub:              pubSub,
		stringInterpreter:   stringPkg.NewInterpreter(),
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
		loopInterpreter:     loop.NewInterpreter(),
	}
}

type _looper struct {
	caller              caller
	loopableInterpreter loopable.Interpreter
	pubSub              pubsub.PubSub
	stringInterpreter   stringPkg.Interpreter
	uniqueStringFactory uniquestring.UniqueStringFactory
	loopInterpreter     loop.Interpreter
}

func (lpr _looper) isLoopEnded(
	index int,
	loop *model.DCGLoop,
) bool {
	if nil != loop.Until && *loop.Until {
		// exit condition provided & met
		return true
	}

	if nil != loop.For {
		if len(loop.For.Each.Array) != 0 {
			return index == len(loop.For.Each.Array)
		}
		if len(loop.For.Each.Object) != 0 {
			return index == len(loop.For.Each.Object)
		}

		// empty array or object
		return true
	}

	return false
}

func (lpr _looper) sortMap(
	m map[string]interface{},
) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}

	sort.Strings(names) //sort keys alphabetically
	return names
}

// @TODO move this to loop interpreter
func (lpr _looper) scopeLoopVars(
	index int,
	outboundScope map[string]*model.Value,
	scgLoop *model.SCGLoop,
	opHandle model.DataHandle,
) error {
	if nil != scgLoop.Index {
		// assign iteration index to requested inboundScope variable
		indexAsFloat64 := float64(index)
		outboundScope[*scgLoop.Index] = &model.Value{
			Number: &indexAsFloat64,
		}
	}
	if nil != scgLoop.For && nil != scgLoop.For.Value {
		var rawValue interface{}

		loopable, err := lpr.loopableInterpreter.Interpret(
			scgLoop.For.Each,
			opHandle,
			outboundScope,
		)

		if nil != loopable.Array {
			rawValue = loopable.Array[index]
		} else {
			sortedNames := lpr.sortMap(loopable.Object)
			name := sortedNames[index]
			rawValue = loopable.Object[name]

			if nil != scgLoop.For.Key {
				// only add key to scope if declared
				outboundScope[*scgLoop.For.Key] = &model.Value{String: &name}
			}
		}

		// interpret value as string since everything is coercible to string
		outboundScope[*scgLoop.For.Value], err = lpr.stringInterpreter.Interpret(
			outboundScope,
			rawValue,
			opHandle,
		)
		if nil != err {
			return err
		}
	}
	return nil
}

func (lpr _looper) Loop(
	ctx context.Context,
	id string,
	inboundScope map[string]*model.Value,
	scg *model.SCG,
	opHandle model.DataHandle,
	parentCallID *string,
	rootOpID string,
) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	outboundScope := map[string]*model.Value{}
	for varName, varData := range inboundScope {
		outboundScope[varName] = varData
	}

	defer func() {
		// defer must be defined before conditional return statements so it always runs
		lpr.pubSub.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				CallEnded: &model.CallEndedEvent{
					CallID:     id,
					RootCallID: rootOpID,
					Outputs:    outboundScope,
				},
			},
		)
	}()

	// store scope shadowed in loop
	shadowedScope := map[string]*model.Value{}
	if nil != scg.Loop.Index {
		shadowedScope[*scg.Loop.Index] = outboundScope[*scg.Loop.Index]
	}
	if nil != scg.Loop.For && nil != scg.Loop.For.Key {
		shadowedScope[*scg.Loop.For.Key] = outboundScope[*scg.Loop.For.Key]
	}
	if nil != scg.Loop.For && nil != scg.Loop.For.Value {
		shadowedScope[*scg.Loop.For.Value] = outboundScope[*scg.Loop.For.Value]
	}

	index := 0
	if err := lpr.scopeLoopVars(
		index,
		outboundScope,
		scg.Loop,
		opHandle,
	); nil != err {
		return err
	}

	// copy scg.Loop & remove from scg since we're already looping
	scgLoop := scg.Loop
	scg.Loop = nil

	// interpret initial iteration of the loop
	dcgLoop, err := lpr.loopInterpreter.Interpret(
		opHandle,
		scgLoop,
		outboundScope,
	)
	if nil != err {
		return err
	}

	for !lpr.isLoopEnded(index, dcgLoop) {
		eventFilterSince := time.Now().UTC()

		callID, err := lpr.uniqueStringFactory.Construct()
		if nil != err {
			return err
		}

		err = lpr.caller.Call(
			ctx,
			callID,
			outboundScope,
			scg,
			opHandle,
			parentCallID,
			rootOpID,
		)
		if nil != err {
			// end looping on any error
			return err
		}

		// subscribe to events
		// @TODO: handle err channel
		eventChannel, _ := lpr.pubSub.Subscribe(
			ctx,
			model.EventFilter{
				Roots: []string{rootOpID},
				Since: &eventFilterSince,
			},
		)

	eventLoop:
		for event := range eventChannel {
			// merge child outboundScope w/ outboundScope, child outboundScope having precedence
			switch {
			case nil != event.OpEnded && event.OpEnded.OpID == callID:
				for name, value := range event.OpEnded.Outputs {
					outboundScope[name] = value
				}
				break eventLoop
			case nil != event.ContainerExited && event.ContainerExited.ContainerID == callID:
				for name, value := range event.ContainerExited.Outputs {
					outboundScope[name] = value
				}
				break eventLoop
			case nil != event.SerialCallEnded && event.SerialCallEnded.CallID == callID:
				for name, value := range event.SerialCallEnded.Outputs {
					outboundScope[name] = value
				}
				break eventLoop
			case nil != event.ParallelCallEnded && event.ParallelCallEnded.CallID == callID:
				break eventLoop
			case nil != event.CallEnded && event.CallEnded.CallID == callID:
				for name, value := range event.CallEnded.Outputs {
					outboundScope[name] = value
				}
				break eventLoop
			}
		}
		cancel()

		index++

		if lpr.isLoopEnded(index, dcgLoop) {
			break
		}

		if err := lpr.scopeLoopVars(
			index,
			outboundScope,
			scgLoop,
			opHandle,
		); nil != err {
			return err
		}

		// interpret next iteration of the loop
		dcgLoop, err = lpr.loopInterpreter.Interpret(
			opHandle,
			scgLoop,
			outboundScope,
		)
		if nil != err {
			return err
		}
	}

	// unshadow shadowed scope
	for name, value := range shadowedScope {
		outboundScope[name] = value
	}
	return nil
}
