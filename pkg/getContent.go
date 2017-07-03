package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
)

func (this _Pkg) GetContent(
	pkgRef string,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {

	handle, err := this.opener.Open(pkgRef)
	if nil != err {
		return nil, err
	}

	return handle.GetContent(contentPath)
}
