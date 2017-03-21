package core

import (
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
	"path/filepath"
)

func (this _core) Create(
	path string,
	description string,
	name string,
) {
	pwd, err := this.vos.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	err = this.pkg.Create(
		pkg.CreateReq{
			Path:        filepath.Join(pwd, path, name),
			Name:        name,
			Description: description,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
