package core

import (
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/sdk-golang/model"
	"path"
)

func (this _core) PkgSetDescription(
	description string,
	pkgRef string,
) {
	if !path.IsAbs(pkgRef) {
		pkgDir := path.Dir(pkgRef)

		if "." == pkgDir {
			// default package location is .opspec subdir of current working directory
			// so if they only provided us a name let's look there
			pkgName := path.Base(pkgRef)
			pkgRef = path.Join(pkgDir, ".opspec", pkgName)
		}

		// make our pkgRef absolute
		pwd, err := this.vos.Getwd()
		if nil != err {
			this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
			return // support fake exiter
		}
		pkgRef = path.Join(pwd, pkgRef)
	}

	err := this.managePackages.SetPackageDescription(
		model.SetPackageDescriptionReq{
			PathToOp:    pkgRef,
			Description: description,
		},
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
