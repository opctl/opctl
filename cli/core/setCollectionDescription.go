package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this _core) SetCollectionDescription(
	description string,
) {
	err := this.bundle.SetCollectionDescription(
		model.SetCollectionDescriptionReq{
			PathToCollection: path.Join(this.workDirPathGetter.Get(), ".opspec"),
			Description:      description,
		},
	)
	if nil != err {
		this.exiter.Exit(ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
