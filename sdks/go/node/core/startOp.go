package core

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
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

	// pull Creds
	if nil != req.Op.PullCreds {
		scgOpCall.PullCreds = &model.SCGPullCreds{
			Username: req.Op.PullCreds.Username,
			Password: req.Op.PullCreds.Password,
		}
	}

	for name := range req.Args {
		// implicitly bind
		scgOpCall.Inputs[name] = ""
	}

	opFile, err := this.opFileGetter.Get(
		ctx,
		*opHandle.Path(),
	)
	if nil != err {
		return "", err
	}
	for name := range opFile.Outputs {
		// implicitly bind
		scgOpCall.Outputs[name] = ""
	}

	dcgOpCall, err := this.opInterpreter.Interpret(
		req.Args,
		scgOpCall,
		opID,
		*opHandle.Path(),
		opID,
	)
	if nil != err {
		return "", err
	}

	go func() {
		defer func() {
			if panicArg := recover(); panicArg != nil {
				// recover from panics; treat as errors
				msg := fmt.Sprint(panicArg, debug.Stack())

				this.pubSub.Publish(
					model.Event{
						Timestamp: time.Now().UTC(),
						OpEnded: &model.OpEnded{
							Error: &model.CallEndedError{
								Message: msg,
							},
							OpID:     dcgOpCall.OpID,
							OpRef:    dcgOpCall.OpPath,
							Outcome:  model.OpOutcomeFailed,
							RootOpID: dcgOpCall.RootOpID,
						},
					},
				)
			}
		}()

		this.opCaller.Call(
			context.Background(),
			dcgOpCall,
			req.Args,
			&opID,
			scgOpCall,
		)
	}()

	return opID, nil

}
