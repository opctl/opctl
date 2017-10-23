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

	// construct scgOpCallInputs
	scgOpCallInputs := map[string]interface{}{}
	for name := range req.Args {
		// map as passed
    scgOpCallInputs[name] = ""
	}

	// construct scgOpCallOutputs
	scgOpCallOutputs := map[string]string{}
	pkgManifest, err := this.pkg.GetManifest(pkgHandle)
	if nil != err {
	  return "", err
  }

  for name := range pkgManifest.Outputs {
    // bind all defined outputs
    scgOpCallOutputs[name] = name
  }

  // construct scgOpCall
  scgOpCall := &model.SCGOpCall{
    Pkg: &model.SCGOpCallPkg{
      Ref: pkgHandle.Ref(),
    },
    Inputs: scgOpCallInputs,
    Outputs: scgOpCallOutputs,
  }

  opId := this.uniqueStringFactory.Construct()

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
