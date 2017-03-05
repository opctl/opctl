package pkg

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func (this pkg) GetOp(
	opPackagePath string,
) (
	opView model.OpView,
	err error,
) {

	return this.opViewFactory.Construct(opPackagePath)

}
