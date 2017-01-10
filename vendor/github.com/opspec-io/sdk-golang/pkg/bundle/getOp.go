package bundle

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func (this _bundle) GetOp(
	opBundlePath string,
) (
	opView model.OpView,
	err error,
) {

	return this.opViewFactory.Construct(opBundlePath)

}
