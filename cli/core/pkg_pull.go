package core

import (
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/pkg"
)

func (this _core) PkgPull(
	pkgRef,
	username,
	password string,
) {

	err := this.pkg.Pull(
		pkgRef,
		&pkg.PullOpts{
			Username: username,
			Password: password,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
