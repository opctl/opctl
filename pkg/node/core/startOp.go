package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func (this _core) StartOp(
	req model.StartOpReq,
) (
	opId string,
	err error,
) {

	opId = this.uniqueStringFactory.Construct()

	// @TODO: remove once caller signature updated to use `inputs chan *variable`
	inputs := make(chan *variable, 150)
	for varName, varValue := range req.Args {
		inputs <- &variable{
			Name:  varName,
			Value: varValue,
		}
	}
	close(inputs)
	go func() {
		this.opCaller.Call(
			inputs,
			make(chan *variable, 150),
			opId,
			req.PkgRef,
			opId,
		)
	}()

	return

}
