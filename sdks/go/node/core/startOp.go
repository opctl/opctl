package core

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/data/git"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

func (this core) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (string, error) {
	if nil == req.Op.PullCreds {
		auth := this.stateStore.TryGetAuth(req.Op.Ref)
		if nil != auth {
			// if explicit pull creds not provided and auth exists; use it
			req.Op.PullCreds = &auth.Creds
		}
	}

	opHandle, err := data.Resolve(
		ctx,
		req.Op.Ref,
		fs.New(),
		git.New(this.dataCachePath, req.Op.PullCreds),
	)
	if nil != err {
		return "", err
	}

	callID, err := this.uniqueStringFactory.Construct()
	if nil != err {
		// end run immediately on any error
		return "", err
	}

	// construct opCallSpec
	opCallSpec := &model.OpCallSpec{
		Ref:     opHandle.Ref(),
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
	}

	// pull Creds
	if nil != req.Op.PullCreds {
		opCallSpec.PullCreds = &model.CredsSpec{
			Username: req.Op.PullCreds.Username,
			Password: req.Op.PullCreds.Password,
		}
	}

	for name := range req.Args {
		// implicitly bind
		opCallSpec.Inputs[name] = ""
	}

	opFile, err := opfile.Get(
		ctx,
		*opHandle.Path(),
	)
	if nil != err {
		return "", err
	}
	for name := range opFile.Outputs {
		// implicitly bind
		opCallSpec.Outputs[name] = ""
	}

	opCtx, cancelOp := context.WithCancel(ctx)
	go func() {
		defer func() {
			if panicArg := recover(); panicArg != nil {
				// recover from panics; treat as errors
				fmt.Println(panicArg, debug.Stack())
			}

			cancelOp()
		}()

		this.caller.Call(
			opCtx,
			callID,
			req.Args,
			&model.CallSpec{
				Op: opCallSpec,
			},
			*opHandle.Path(),
			nil,
			callID,
		)
	}()

	return callID, nil
}
