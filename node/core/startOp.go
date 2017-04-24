package core

import (
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

func (this _core) StartOp(
	req model.StartOpReq,
) (
	opId string,
	err error,
) {

	pkgBasePath := path.Dir(req.PkgRef)
	pkgName := path.Base(req.PkgRef)

	opId = this.uniqueStringFactory.Construct()

	// construct scgOpCall
	scgOpCall := &model.SCGOpCall{
		Pkg: &model.SCGOpCallPkg{
			Ref: pkgName,
		},
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
			pkgBasePath,
			opId,
			scgOpCall,
		)
	}()

	return

}
