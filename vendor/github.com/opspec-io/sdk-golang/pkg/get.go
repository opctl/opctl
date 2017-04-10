package pkg

import (
	"github.com/opspec-io/sdk-golang/model"
)

// Get gets a package according to opspec package resolution rules
func (this pkg) Get(
	req *GetReq,
) (*model.PkgManifest, error) {
	return this.getter.Get(req)
}
