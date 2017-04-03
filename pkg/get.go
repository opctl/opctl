package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
)

func (this pkg) Get(
	req *GetReq,
) (
	packageView *model.PackageView,
	err error,
) {

	packageView, err = this.getter.Get(req)

	return
}
