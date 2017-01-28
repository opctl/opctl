package core

//go:generate counterfeiter -o ./fakeSerialCaller.go --fake-name fakeSerialCaller ./ serialCaller

import (
	"fmt"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type serialCaller interface {
	Call(
		parentScope map[string]*model.Data,
		opGraphId string,
		opRef string,
		serialCall []*model.Scg,
	) (
		err error,
	)
}

func newSerialCaller(
	caller caller,
	uniqueStringFactory uniquestring.UniqueStringFactory,
) serialCaller {

	return _serialCaller{
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _serialCaller struct {
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _serialCaller) Call(
	parentScope map[string]*model.Data,
	opGraphId string,
	opRef string,
	serialCall []*model.Scg,
) (
	err error,
) {
	currentScope := map[string]*model.Data{}
	for varName, varData := range parentScope {
		currentScope[varName] = varData
	}

	var childScope map[string]*model.Data
	for _, call := range serialCall {
		fmt.Printf("serialCaller.scope:\n %#v\n", currentScope)
		childScope, err = this.caller.Call(
			this.uniqueStringFactory.Construct(),
			currentScope,
			call,
			opRef,
			opGraphId,
		)
		if nil != err {
			// end run immediately on any error
			return
		}

		// apply outputs to current scope
		for varName, varData := range childScope {
			currentScope[varName] = varData
		}
	}

	return

}
