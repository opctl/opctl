package core

import (
	"context"
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _core) StartOp(
	ctx context.Context,
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

	opHandle, err := this.data.Resolve(
		ctx,
		req.Pkg.Ref,
		this.data.NewFSProvider(),
		this.data.NewGitProvider(this.dataCachePath, pullCreds),
	)
	if nil != err {
		return "", err
	}

	opID, err := this.uniqueStringFactory.Construct()
	if nil != err {
		// end run immediately on any error
		return "", err
	}

	// construct scgOpCall
	scgOpCall := &model.SCGOpCall{
		Pkg: &model.SCGOpCallPkg{
			Ref: opHandle.Ref(),
		},
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
	}
	for name := range req.Args {
		// implicitly bind
		scgOpCall.Inputs[name] = ""
	}

	opDotYml, err := this.dotYmlGetter.Get(
		context.TODO(),
		opHandle,
	)
	if nil != err {
		return "", err
	}
	for name := range opDotYml.Outputs {
		// implicitly bind
		scgOpCall.Outputs[name] = ""
	}

	go func() {
		this.opCaller.Call(
			req.Args,
			opID,
			opHandle,
			opID,
			scgOpCall,
		)
	}()

	return opID, nil

}
