package core

import (
	"bytes"
	"context"
	"fmt"

	"github.com/opctl/opctl/cli/util/cliexiter"
)

func (this _core) OpValidate(
	ctx context.Context,
	opRef string,
) {
	opDirHandle := this.dataResolver.Resolve(
		opRef,
		nil,
	)

	errs := this.opValidator.Validate(
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
		this.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf(
				`%v
-`, messageBuffer.String()),
			Code: 1})
	} else {
		this.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf("%v is valid", opDirHandle.Ref()),
		})
	}
}
