package core

import (
	"context"
	"fmt"
	"github.com/opctl/opctl/util/cliexiter"
	"text/tabwriter"
)

func (this _core) Ls(
	ctx context.Context,
	dirRef string,
) {
	_tabWriter := new(tabwriter.Writer)
	defer _tabWriter.Flush()
	_tabWriter.Init(this.writer, 0, 8, 1, '\t', 0)

	fmt.Fprintln(_tabWriter, "NAME\tVERSION\tDESCRIPTION")

	dirHandle := this.dataResolver.Resolve(
		dirRef,
		nil,
	)

	packages, err := this.opLister.List(
		ctx,
		dirHandle,
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
