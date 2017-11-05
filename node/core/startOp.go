package core

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _core) StartOp(
	req model.StartOpReq,
) (string, error) {
	if nil == req.Pkg {
		return "", errors.New("pkg required")
	}

	var pullCreds *model.PullCreds
	if nil != req.Pkg.PullCreds {
		pullCreds = &model.PullCreds{
			Username: req.Pkg.PullCreds.Username,
			Password: req.Pkg.PullCreds.Password,
		}
	}

	pkgHandle, err := this.pkg.Resolve(
		req.Pkg.Ref,
		this.pkg.NewFSProvider(),
		this.pkg.NewGitProvider(this.pkgCachePath, pullCreds),
	)
	if nil != err {
		return "", err
	}

	opId := this.uniqueStringFactory.Construct()

	// construct scgOpCall
	scgOpCall := &model.SCGOpCall{
		Pkg: &model.SCGOpCallPkg{
			Ref: pkgHandle.Ref(),
		},
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
	}
	for name := range req.Args {
		// implicitly bind
		scgOpCall.Inputs[name] = ""
	}

	pkgManifest, err := this.pkg.GetManifest(pkgHandle)
	if nil != err {
		return "", err
	}
	for name := range pkgManifest.Outputs {
		// implicitly bind
		scgOpCall.Outputs[name] = ""
	}

	go func() {
		this.opCaller.Call(
			req.Args,
			opId,
			pkgHandle,
			opId,
			scgOpCall,
		)
	}()

	return opId, nil

}
