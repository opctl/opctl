package core

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (this _core) StartOp(
	req model.StartOpReq,
) (string, error) {
	if nil == req.Pkg {
		return "", errors.New("pkg required")
	}

	pkgBasePath := filepath.Dir(req.Pkg.Ref)
	pkgName := filepath.Base(req.Pkg.Ref)

	opId := this.uniqueStringFactory.Construct()

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

	return opId, nil

}
