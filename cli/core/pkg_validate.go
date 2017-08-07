package core

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
)

func (this _core) PkgValidate(
	pkgRef string,
) {
	pkgHandle := this.pkgResolver.Resolve(
		pkgRef,
		nil,
	)

	errs := this.pkg.Validate(pkgHandle)
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
			Message: fmt.Sprintf("%v is valid", pkgHandle.Ref()),
		})
	}
}
