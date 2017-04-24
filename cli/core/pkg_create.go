package core

import (
	"github.com/opctl/opctl/util/cliexiter"
	pathPkg "path"
)

func (this _core) Create(
	path string,
	description string,
	name string,
) {
	cwd, err := this.vos.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	err = this.pkg.Create(
		pathPkg.Join(cwd, path, name),
		name,
		description,
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
