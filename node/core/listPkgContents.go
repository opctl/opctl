package core

import "github.com/opspec-io/sdk-golang/model"

func (this _core) ListPkgContents(
	pkgRef string,
) (
	[]*model.PkgContent,
	error,
) {
	return this.pkg.ListContents(pkgRef)
}
