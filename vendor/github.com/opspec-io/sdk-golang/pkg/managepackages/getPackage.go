package managepackages

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func (this managePackages) GetPackage(
	packageRef string,
) (
	packageView model.PackageView,
	err error,
) {

	return this.packageViewFactory.Construct(packageRef)

}
