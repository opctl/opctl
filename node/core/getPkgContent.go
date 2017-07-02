package core

import "github.com/opspec-io/sdk-golang/model"

func (this _core) GetPkgContent(
	pkgRef,
	contentPath string,
) (model.ReadSeekCloser, error) {
	return this.pkg.GetContent(pkgRef, contentPath)
}
