package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this _core) CreateOp(
	collection string,
	description string,
	name string,
) {
	err := this.bundle.CreateOp(
		model.CreateOpReq{
			Path:        path.Join(this.workDirPathGetter.Get(), collection, name),
			Name:        name,
			Description: description,
		},
	)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
