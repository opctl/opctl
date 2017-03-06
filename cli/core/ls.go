package core

import (
	"fmt"
	"github.com/opspec-io/opctl/util/cliexiter"
	"path/filepath"
	"text/tabwriter"
)

func (this _core) ListOpsInCollection(
	collection string,
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

	ops, err := this.pkg.GetCollection(
		filepath.Join(pwd, collection),
	)
	if nil != err {
		this.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for _, op := range ops.Ops {

		fmt.Fprintf(_tabWriter, "%v\t%v", op.Name, op.Description)
		fmt.Fprintln(_tabWriter)

	}
}
