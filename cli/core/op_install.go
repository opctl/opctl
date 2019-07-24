package core

import (
	"context"
	"github.com/opctl/opctl/cli/util/cliexiter"
	"github.com/opctl/opctl/sdks/go/types"
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
		&types.PullCreds{
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
