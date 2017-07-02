package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
	"path/filepath"
)

func (this _Pkg) GetContent(
	pkgRef string,
	contentPath string,
) (
	model.ReadSeekCloser,
	error,
) {
	return this.os.Open(filepath.Join(pkgRef, contentPath))
}
