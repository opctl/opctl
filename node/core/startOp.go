package core

import (
	"github.com/opspec-io/sdk-golang/model"
)

func (this _core) StartOp(
	req model.StartOpReq,
) (
	opId string,
	err error,
) {

	opId = this.uniqueStringFactory.Construct()

	// construct scgOpCall
	scgOpCall := &model.ScgOpCall{
		Ref:    req.PkgRef,
		Inputs: map[string]string{},
	}
	for name := range req.Args {
		// map as passed
		scgOpCall.Inputs[name] = name
	}

	go func() {
		this.opCaller.Call(
			req.Args,
			opId,
			req.PkgRef,
			opId,
			scgOpCall,
		)
	}()

	return

}
