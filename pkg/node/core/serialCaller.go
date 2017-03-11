package core

//go:generate counterfeiter -o ./fakeSerialCaller.go --fake-name fakeSerialCaller ./ serialCaller

import (
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type serialCaller interface {
	// Executes a serial call
	Call(
		inputs chan *variable,
		outputs chan *variable,
		rootOpId string,
		pkgRef string,
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
	inputs chan *variable,
	outputs chan *variable,
	rootOpId string,
	pkgRef string,
	scgSerialCall []*model.Scg,
) (
	err error,
) {
	scope := map[string]*model.Data{}
	for input := range inputs {
		scope[input.Name] = input.Value
	}

	for _, scgCall := range scgSerialCall {

		subInputs := make(chan *variable, 150)
		for varName, varValue := range scope {
			subInputs <- &variable{
				Name:  varName,
				Value: varValue,
			}
		}
		close(subInputs)

		subOutputs := make(chan *variable, 150)

		err = this.caller.Call(
			this.uniqueStringFactory.Construct(),
			subInputs,
			subOutputs,
			scgCall,
			pkgRef,
			rootOpId,
		)
		if nil != err {
			// end immediately on any error
			return
		}

		if scgOpCall := scgCall.Op; nil != scgCall.Op {
			// apply bound child outputs to current scope
			for subOutput := range subOutputs {
				for currentScopeVarName, childScopeVarName := range scgOpCall.Outputs {
					if currentScopeVarName == subOutput.Name || childScopeVarName == subOutput.Name {
						scope[currentScopeVarName] = subOutput.Value
					}
				}
			}
		} else {
			// apply child outputs to current scope
			for subOutput := range subOutputs {
				scope[subOutput.Name] = subOutput.Value
			}
		}
	}

	// @TODO: stream outputs from last child
	for varName, varValue := range scope {
		outputs <- &variable{Name: varName, Value: varValue}
	}
	close(outputs)

	return

}
