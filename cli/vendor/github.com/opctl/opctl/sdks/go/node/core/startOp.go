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

	go func() {
		defer func() {
			if panicArg := recover(); panicArg != nil {
				// recover from panics; treat as errors
				msg := fmt.Sprint(panicArg, debug.Stack())

				this.pubSub.Publish(
					model.Event{
						Timestamp: time.Now().UTC(),
						CallEnded: &model.CallEnded{
							Error: &model.CallEndedError{
								Message: msg,
							},
							CallID:     callID,
							Ref:        *opHandle.Path(),
							Outcome:    model.OpOutcomeFailed,
							RootCallID: callID,
						},
					},
				)
			}
		}()

		this.caller.Call(
			context.Background(),
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
