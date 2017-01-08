package core

import (
	"fmt"
	"github.com/opspec-io/opctl/util/uniquestring"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

type serialCaller interface {
	Call(
		inputs map[string]*model.Data,
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

	return &_serialCaller{
		caller:              caller,
		uniqueStringFactory: uniqueStringFactory,
	}

}

type _serialCaller struct {
	caller              caller
	uniqueStringFactory uniquestring.UniqueStringFactory
}

func (this _serialCaller) Call(
	inputs map[string]*model.Data,
	opGraphId string,
	opRef string,
	serialCall []*model.Scg,
) (
	err error,
) {
	// construct scope
	// Why not just use inputs directly? maps passed by ref in go.. mutating parent scope would be invalid
	scope := map[string]*model.Data{}
	for inputName, inputValue := range inputs {
		scope[inputName] = inputValue
	}

	var outputs map[string]*model.Data
	for _, call := range serialCall {
		fmt.Printf("serialCaller.scope:\n %#v\n", scope)
		outputs, err = this.caller.Call(
			this.uniqueStringFactory.Construct(),
			scope,
			call,
			opRef,
			opGraphId,
		)
		if nil != err {
			// end run immediately on any error
			return
		}

		// apply outputs to current scope
		for outputName, outputValue := range outputs {
			scope[outputName] = outputValue
		}
	}

	return

}
