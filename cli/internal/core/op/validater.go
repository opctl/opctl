package op

import (
	"context"
	"fmt"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// Validater exposes the "op validate" sub command
type Validater interface {
	Validate(
		ctx context.Context,
		opRef string,
	)
}

// newValidater returns an initialized "op validate" sub command
func newValidater(
	cliExiter cliexiter.CliExiter,
	dataResolver dataresolver.DataResolver,
) Validater {
	return _validater{
		cliExiter:    cliExiter,
		dataResolver: dataResolver,
	}
}

type _validater struct {
	cliExiter    cliexiter.CliExiter
	dataResolver dataresolver.DataResolver
}

func (ivkr _validater) Validate(
	ctx context.Context,
	opRef string,
) {
	opDirHandle, err := ivkr.dataResolver.Resolve(
		opRef,
		nil,
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: err.Error(),
			Code:    1},
		)
	}

	if err := opspec.Validate(
		ctx,
		*opDirHandle.Path(),
	); nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: err.Error(),
			Code:    1},
		)
	} else {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf("%v is valid", opDirHandle.Ref()),
		})
	}
}
