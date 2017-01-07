package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this _core) CreateCollection(
	description string,
	name string,
) {
	err := this.bundle.CreateCollection(
		model.CreateCollectionReq{
			Path:        path.Join(this.workDirPathGetter.Get(), name),
			Name:        name,
			Description: description,
		},
	)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
