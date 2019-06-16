package core

import (
	"context"
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/util/cliexiter"
)

func (this _core) OpInstall(
	ctx context.Context,
	path,
	opRef,
	username,
	password string,
) {

	opDirHandle := this.dataResolver.Resolve(
		opRef,
		&model.PullCreds{
			Username: username,
			Password: password,
		},
	)

	if err := this.opInstaller.Install(
		ctx,
		path,
		opDirHandle,
	); nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
	}

}
