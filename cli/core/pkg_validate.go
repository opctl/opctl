package core

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
	"path"
)

func (this _core) PkgValidate(
	pkgRef string,
) {
	originalPkgRef := pkgRef
	if !path.IsAbs(pkgRef) {
		pkgDir := path.Dir(pkgRef)

		if "." == pkgDir {
			// default package location is .opspec subdir of current working directory
			// so if they only provided us a name let's look there
			pkgName := path.Base(pkgRef)
			pkgRef = path.Join(pkgDir, ".opspec", pkgName)
		}

		// make our pkgRef absolute
		cwd, err := this.os.Getwd()
		if nil != err {
			this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
			return // support fake exiter
		}
		pkgRef = path.Join(cwd, pkgRef)
	}

	errs := this.pkg.Validate(pkgRef)
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
			Message: fmt.Sprintf("%v is valid", originalPkgRef),
		})
	}
}
