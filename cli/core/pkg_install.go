package core

import (
	"context"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
)

func (this _core) PkgInstall(
	path,
	pkgRef,
	username,
	password string,
) {

	pkgHandle := this.pkgResolver.Resolve(
		pkgRef,
		&model.PullCreds{
			Username: username,
			Password: password,
		},
	)

	if err := this.pkg.Install(context.TODO(), path, pkgHandle); nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
	}

}
