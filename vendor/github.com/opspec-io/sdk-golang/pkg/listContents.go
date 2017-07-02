package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
)

func (this _Pkg) ListContents(
	pkgRef string,
) (
	[]*model.PkgContent,
	error,
) {
	handle, err := this.opener.Open(pkgRef)
	if nil != err {
		return nil, err
	}

	return handle.ListContents()
}
