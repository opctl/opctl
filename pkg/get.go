package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
)

// Get gets a package according to opspec package resolution rules
func (this _Pkg) Get(
	pkgRef string,
) (*model.PkgManifest, error) {
	return this.manifestUnmarshaller.Unmarshal(pkgRef)
}
