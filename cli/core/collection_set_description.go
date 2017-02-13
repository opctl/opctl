package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this _core) SetCollectionDescription(
	description string,
) {
	pwd, err := this.vos.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	err = this.bundle.SetCollectionDescription(
		model.SetCollectionDescriptionReq{
			PathToCollection: path.Join(pwd, ".opspec"),
			Description:      description,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
