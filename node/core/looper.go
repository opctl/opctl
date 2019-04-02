package core

//go:generate counterfeiter -o ./fakeLooper.go --fake-name fakeLooper ./ looper

import (
	"github.com/opctl/sdk-golang/opspec/interpreter/array"

	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop"
	stringPkg "github.com/opctl/sdk-golang/opspec/interpreter/string"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

type looper interface {
	// Loop loops a call
	Loop(
		id string,
		scope map[string]*model.Value,
		scg *model.SCG,
		opHandle model.DataHandle,
		rootOpID string,
	) error
}

func newLooper(
	caller caller,
	pubSub pubsub.PubSub,
) looper {
	return _looper{
		arrayInterpreter:    array.NewInterpreter(),
		caller:              caller,
		pubSub:              pubSub,
		stringInterpreter:   stringPkg.NewInterpreter(),
		uniqueStringFactory: uniquestring.NewUniqueStringFactory(),
		loopInterpreter:     loop.NewInterpreter(),
	}
}

type _looper struct {
	arrayInterpreter    array.Interpreter
	caller              caller
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

	if nil != loop.For && index == len(loop.For.Each.Array) {
		// each provided & iteration complete
		return true
	}

	return false
}

func (lpr _looper) scopeLoopVars(
	index int,
	scope map[string]*model.Value,
	scgLoop *model.SCGLoop,
	opHandle model.DataHandle,
) error {
	if nil != scgLoop.Index {
		// assign iteration index to requested scope variable
		indexAsFloat64 := float64(index)
		scope[*scgLoop.Index] = &model.Value{
			Number: &indexAsFloat64,
		}
	}
	if nil != scgLoop.For && nil != scgLoop.For.Value {
		each, err := lpr.arrayInterpreter.Interpret(
			scope,
			scgLoop.For.Each,
			opHandle,
		)

		// interpret value as string since everything is coercible to string
		scope[*scgLoop.For.Value], err = lpr.stringInterpreter.Interpret(
			scope,
			each.Array[index],
			opHandle,
		)
		if nil != err {
			return err
		}
	}
	return nil
}

func (lpr _looper) Loop(
	id string,
	scope map[string]*model.Value,
	scg *model.SCG,
	opHandle model.DataHandle,
	rootOpID string,
) error {
	// store scope shadowed in loop
	shadowedScope := map[string]*model.Value{}
	if nil != scg.Loop.Index {
		shadowedScope[*scg.Loop.Index] = scope[*scg.Loop.Index]
	}
	if nil != scg.Loop.For && nil != scg.Loop.For.Value {
		shadowedScope[*scg.Loop.For.Value] = scope[*scg.Loop.For.Value]
	}

	index := 0
	if err := lpr.scopeLoopVars(
		index,
		scope,
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
		scope,
	)
	if nil != err {
		return err
	}

	if lpr.isLoopEnded(index, dcgLoop) {
		return nil
	}

	for {
		callID, err := lpr.uniqueStringFactory.Construct()
		if nil != err {
			return err
		}

		err = lpr.caller.Call(
			callID,
			scope,
			scg,
			opHandle,
			rootOpID,
		)
		if nil != err {
			// end looping on any error
			return err
		}

		index++

		if lpr.isLoopEnded(index, dcgLoop) {
			break
		}

		if err := lpr.scopeLoopVars(
			index,
			scope,
			scgLoop,
			opHandle,
		); nil != err {
			return err
		}

		// interpret next iteration of the loop
		dcgLoop, err = lpr.loopInterpreter.Interpret(
			opHandle,
			scgLoop,
			scope,
		)
		if nil != err {
			return err
		}
	}

	// unshadow shadowed scope
	for name, value := range shadowedScope {
		scope[name] = value
	}
	return nil
}
