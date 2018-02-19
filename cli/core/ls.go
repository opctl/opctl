package core

import (
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
	"path/filepath"
	"text/tabwriter"
)

func (this _core) PkgLs(
	path string,
) {
	_tabWriter := new(tabwriter.Writer)
	defer _tabWriter.Flush()
	_tabWriter.Init(this.writer, 0, 8, 1, '\t', 0)

	fmt.Fprintln(_tabWriter, "NAME\tVERSION\tDESCRIPTION")

	cwd, err := this.os.Getwd()
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	packages, err := this.pkg.ListOps(
		filepath.Join(cwd, path),
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for _, _package := range packages {

		fmt.Fprintf(_tabWriter, "%v\t%v\t%v", _package.Name, _package.Version, _package.Description)
		fmt.Fprintln(_tabWriter)

	}
}
