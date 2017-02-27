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
		outboundScope map[string]*model.Data,
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
	outboundScope map[string]*model.Data,
	err error,
) {
	outboundScope = map[string]*model.Data{}
	for varName, varData := range inboundScope {
		outboundScope[varName] = varData
	}

	for _, scgCall := range scgSerialCall {
		var childOutboundScope map[string]*model.Data
		childOutboundScope, err = this.caller.Call(
			this.uniqueStringFactory.Construct(),
			outboundScope,
			scgCall,
			opRef,
			opGraphId,
		)
		if nil != err {
			// end run immediately on any error
			return
		}

		if scgOpCall := scgCall.Op; nil != scgCall.Op {
			// apply bound child outputs to current scope
			for currentScopeVarName, childScopeVarName := range scgOpCall.Outputs {
				if "" == childScopeVarName {
					// if no explicit childScopeVarName provided; use currentScopeVarName (assume same)
					childScopeVarName = currentScopeVarName
				}
				outboundScope[currentScopeVarName] = childOutboundScope[childScopeVarName]
			}
		} else {
			// apply child outputs to current scope
			for varName, varData := range childOutboundScope {
				outboundScope[varName] = varData
			}
		}
	}

	return

}
