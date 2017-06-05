package core

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
)

func (this _core) PkgValidate(
	pkgRef string,
) {
	cwd, err := this.os.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	parsedPkgRef, err := this.pkg.ParseRef(pkgRef)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	pkgPath, ok := this.pkg.Resolve(parsedPkgRef, cwd)
	if !ok {
		msg := fmt.Sprintf("Unable to resolve package '%v' from '%v'", pkgRef, cwd)
		this.cliExiter.Exit(cliexiter.ExitReq{Message: msg, Code: 1})
		return // support fake exiter
	}

	errs := this.pkg.Validate(pkgPath)
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
			Message: fmt.Sprintf("%v is valid", pkgPath),
		})
	}
}
