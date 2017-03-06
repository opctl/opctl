package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path"
)

func (this _core) SetOpDescription(
	collection string,
	description string,
	name string,
) {
	pwd, err := this.vos.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	err = this.pkg.SetOpDescription(
		model.SetOpDescriptionReq{
			PathToOp:    path.Join(pwd, collection, name),
			Description: description,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
