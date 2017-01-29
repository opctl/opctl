package core

//go:generate counterfeiter -o ./fakeSerialCaller.go --fake-name fakeSerialCaller ./ serialCaller

import (
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type serialCaller interface {
	// Executes a serial call
	Call(
		inboundScope map[string]*model.Data,
		opGraphId string,
		opRef string,
		scgSerialCall []*model.Scg,
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
	inboundScope map[string]*model.Data,
	opGraphId string,
	opRef string,
	scgSerialCall []*model.Scg,
) (
	err error,
) {
	currentScope := map[string]*model.Data{}
	for varName, varData := range inboundScope {
		currentScope[varName] = varData
	}

	var childOutboundScope map[string]*model.Data
	for _, call := range scgSerialCall {
		childOutboundScope, err = this.caller.Call(
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
		for varName, varData := range childOutboundScope {
			currentScope[varName] = varData
		}
	}

	return

}
