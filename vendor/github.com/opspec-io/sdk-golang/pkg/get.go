package pkg

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

func (this pkg) Get(
	pkgRef string,
) (
	packageView model.PackageView,
	err error,
) {

	// 1) ensure valid
	errs := this.validator.Validate(pkgRef)
	if len(errs) > 0 {
		messageBuffer := bytes.NewBufferString(
			fmt.Sprint(`
-
  Error(s):`))
		for _, validationError := range errs {
			messageBuffer.WriteString(fmt.Sprintf(`
    - %v`, validationError.Error()))
		}
		err = fmt.Errorf(
			`%v
-`, messageBuffer.String())
	}
	if nil != err {
		return
	}

	// 2) construct view
	packageView, err = this.packageViewFactory.Construct(pkgRef)

	return
}
