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

	go func() {
		this.opCaller.Call(
			req.Args,
			opId,
			req.PkgRef,
			opId,
		)
	}()

	return

}
