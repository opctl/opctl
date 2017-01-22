package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path/filepath"
)

func (this _core) CreateCollection(
	description string,
	name string,
) {
	pwd, err := this.vos.Getwd()
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	err = this.bundle.CreateCollection(
		model.CreateCollectionReq{
			Path:        filepath.Join(pwd, name),
			Name:        name,
			Description: description,
		},
	)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
