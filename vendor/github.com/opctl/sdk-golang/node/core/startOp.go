package core

import (
	"context"
	"github.com/opctl/sdk-golang/model"
)

func (this _core) StartOp(
	ctx context.Context,
	req model.StartOpReq,
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
	scgOpCall := &model.SCGOpCall{
		Ref:     opHandle.Ref(),
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
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

	dcg, err := this.opInterpreter.Interpret(
		req.Args,
		scgOpCall,
		opID,
		opHandle,
		opID,
	)
	if nil != err {
		return "", err
	}

	go func() {
		this.opCaller.Call(
			dcg,
			req.Args,
			scgOpCall,
		)
	}()

	return opID, nil

}
