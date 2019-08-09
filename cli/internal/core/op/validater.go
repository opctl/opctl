package op

import (
	"bytes"
	"context"
	"fmt"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	opspec "github.com/opctl/opctl/sdks/go/opspec"
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
		opValidator:  opspec.NewValidator(),
	}
}

type _validater struct {
	cliExiter    cliexiter.CliExiter
	dataResolver dataresolver.DataResolver
	opValidator  opspec.Validator
}

func (ivkr _validater) Validate(
	ctx context.Context,
	opRef string,
) {
	opDirHandle := ivkr.dataResolver.Resolve(
		opRef,
		nil,
	)

	errs := ivkr.opValidator.Validate(
		ctx,
		opDirHandle,
	)
	if len(errs) > 0 {
		messageBuffer := bytes.NewBufferString(
			fmt.Sprint(`
-
  Error(s):`))
		for _, validationError := range errs {
			messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
		}
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf(
				`%v
-`, messageBuffer.String()),
			Code: 1})
	} else {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf("%v is valid", opDirHandle.Ref()),
		})
	}
}
