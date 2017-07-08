package core

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
)

func (this _core) PkgValidate(
	pkgRef string,
) {
	cwd, err := this.os.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	pkgHandle, err := this.pkg.Resolve(pkgRef, &pkg.ResolveOpts{BasePath: cwd})
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{
			Message: fmt.Sprintf("Unable to resolve package '%v' from '%v'; error was: %v", pkgRef, cwd, err),
			Code:    1})
		return // support fake exiter
	}

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
