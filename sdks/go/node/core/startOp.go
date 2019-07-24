package core

import (
	"context"
	"github.com/opctl/opctl/sdks/go/types"
)

func (this _core) StartOp(
	ctx context.Context,
	req types.StartOpReq,
) (string, error) {
	opHandle, err := this.data.Resolve(
		ctx,
		req.Op.Ref,
		this.data.NewFSProvider(),
		this.data.NewGitProvider(this.dataCachePath, req.Op.PullCreds),
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
	scgOpCall := &types.SCGOpCall{
		Ref:     opHandle.Ref(),
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
	}

	// pull Creds
	if nil != req.Op.PullCreds {
		scgOpCall.PullCreds = &types.SCGPullCreds{
			Username: req.Op.PullCreds.Username,
			Password: req.Op.PullCreds.Password,
		}
	}

	for name := range req.Args {
		// implicitly bind
		scgOpCall.Inputs[name] = ""
	}

	opDotYml, err := this.dotYmlGetter.Get(
		ctx,
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
		this.caller.Call(
			// call in background context
			context.Background(),
			opID,
			req.Args,
			&types.SCG{
				Op: scgOpCall,
			},
			opHandle,
			nil,
			opID,
		)
	}()

	return opID, nil

}
