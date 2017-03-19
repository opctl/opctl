package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
)

func (this pkg) Get(
	pkgRef string,
) (
	packageView model.PackageView,
	err error,
) {

	return this.packageViewFactory.Construct(pkgRef)

}
