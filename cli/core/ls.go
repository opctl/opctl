package core

import (
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
	pathPkg "path"
	"text/tabwriter"
)

func (this _core) ListPackages(
	path string,
) {
	_tabWriter := new(tabwriter.Writer)
	defer _tabWriter.Flush()
	_tabWriter.Init(this.writer, 0, 8, 1, '\t', 0)

	fmt.Fprintln(_tabWriter, "NAME\tDESCRIPTION")

	pwd, err := this.vos.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	packages, err := this.pkg.ListPackagesInDir(
		pathPkg.Join(pwd, path),
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for _, _package := range packages {

		fmt.Fprintf(_tabWriter, "%v\t%v", _package.Name, _package.Description)
		fmt.Fprintln(_tabWriter)

	}
}
