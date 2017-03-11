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
	inboundScope map[string]*model.Data,
	outputs chan *variable,
	rootOpId string,
	pkgRef string,
	scgSerialCall []*model.Scg,
) (
	err error,
) {
	outboundScope := map[string]*model.Data{}
	for varName, varData := range inboundScope {
		outboundScope[varName] = varData
	}

	for _, scgCall := range scgSerialCall {
		subOutputs := make(chan *variable, 150)
		err = this.caller.Call(
			this.uniqueStringFactory.Construct(),
			outboundScope,
			subOutputs,
			scgCall,
			pkgRef,
			rootOpId,
		)
		if nil != err {
			// end run immediately on any error
			return
		}

		if scgOpCall := scgCall.Op; nil != scgCall.Op {
			// apply bound child outputs to current scope
			for subOutput := range subOutputs {
				for currentScopeVarName, childScopeVarName := range scgOpCall.Outputs {
					if currentScopeVarName == subOutput.Name || childScopeVarName == subOutput.Name {
						outboundScope[currentScopeVarName] = subOutput.Value
					}
				}
			}
		} else {
			// apply child outputs to current scope
			for subOutput := range subOutputs {
				outboundScope[subOutput.Name] = subOutput.Value
			}
		}
	}

	// @TODO: stream outputs from last child
	for varName, varValue := range outboundScope {
		outputs <- &variable{Name: varName, Value: varValue}
	}
	close(outputs)

	return

}
