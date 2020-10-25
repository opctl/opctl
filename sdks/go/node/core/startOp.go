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

	// construct opCallSpec
	opCallSpec := &model.OpCallSpec{
		Ref:     opHandle.Ref(),
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
	}

	// pull Creds
	if nil != req.Op.PullCreds {
		opCallSpec.PullCreds = &model.PullCredsSpec{
			Username: req.Op.PullCreds.Username,
			Password: req.Op.PullCreds.Password,
		}
	}

	for name := range req.Args {
		// implicitly bind
		opCallSpec.Inputs[name] = ""
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
		opCallSpec.Outputs[name] = ""
	}

	opCall, err := this.opInterpreter.Interpret(
		req.Args,
		opCallSpec,
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
							OpID:     opCall.OpID,
							OpRef:    opCall.OpPath,
							Outcome:  model.OpOutcomeFailed,
							RootOpID: opCall.RootOpID,
						},
					},
				)
			}
		}()

		this.opCaller.Call(
			context.Background(),
			opCall,
			req.Args,
			&opID,
			opCallSpec,
		)
	}()

	return opID, nil

}
