package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"path/filepath"
)

func (this _core) CreateOp(
	collection string,
	description string,
	name string,
) {
	pwd, err := this.vos.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	err = this.bundle.CreateOp(
		model.CreateOpReq{
			Path:        filepath.Join(pwd, collection, name),
			Name:        name,
			Description: description,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
