package node

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime/debug"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

func (this core) StartOp(
	ctx context.Context,
	req model.StartOpReq,
) (string, error) {
	opDir, err := this.ResolveData(
		ctx,
		req.Op.Ref,
		req.Op.PullCreds,
	)
	if err != nil {
		return "", err
	}

	callID, err := uniquestring.Construct()
	if err != nil {
		// end run immediately on any error
		return "", err
	}

	// construct opCallSpec
	opCallSpec := &model.OpCallSpec{
		Ref:     opDir.Ref(),
		Inputs:  map[string]interface{}{},
		Outputs: map[string]string{},
	}

	// pull Creds
	if req.Op.PullCreds != nil {
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
		opDir,
	)
	if err != nil {
		return "", err
	}
	for name := range opFile.Outputs {
		// implicitly bind
		opCallSpec.Outputs[name] = ""
	}

	opCtx, cancelOp := context.WithCancel(ctx)
	go func() {
		defer func() {
			if panic := recover(); panic != nil {
				// recover from panics; treat as errors
				fmt.Printf("recovered from panic: %s\n%s\n", panic, string(debug.Stack()))
			}

			cancelOp()
		}()

		// this relies on op existing locally which is currently always true since we only use fs & git as data providers
		opPath := opDir.Ref()
		if !filepath.IsAbs(opPath) {
			opPath = filepath.Join(this.dataCachePath, "ops", opPath)
		}

		this.caller.Call(
			opCtx,
			callID,
			req.Args,
			&model.CallSpec{
				Op: opCallSpec,
			},
			opPath,
			nil,
			callID,
		)
	}()

	return callID, nil
}
