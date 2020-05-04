package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	op "github.com/opctl/opctl/sdks/go/opspec"
)

// Lser exposes the "ls" command
type Lser interface {
	Ls(
		ctx context.Context,
		dirRef string,
	)
}

// newLser returns an initialized "ls" command
func newLser(
	cliExiter cliexiter.CliExiter,
	cliOutput clioutput.CliOutput,
	dataResolver dataresolver.DataResolver,
) Lser {
	return _lsInvoker{
		cliExiter:    cliExiter,
		cliOutput:    cliOutput,
		dataResolver: dataResolver,
		opLister:     op.NewLister(),
		writer:       os.Stdout,
	}
}

type _lsInvoker struct {
	cliExiter    cliexiter.CliExiter
	cliOutput    clioutput.CliOutput
	dataResolver dataresolver.DataResolver
	opLister     op.Lister
	writer       io.Writer
}

func (ivkr _lsInvoker) Ls(
	ctx context.Context,
	dirRef string,
) {
	_tabWriter := new(tabwriter.Writer)
	defer _tabWriter.Flush()
	_tabWriter.Init(ivkr.writer, 0, 8, 1, '\t', 0)

	fmt.Fprintln(_tabWriter, "REF\tDESCRIPTION")

	dirHandle := ivkr.dataResolver.Resolve(
		dirRef,
		nil,
	)

	opsByPath, err := ivkr.opLister.List(
		ctx,
		dirHandle,
	)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	cwd, err := os.Getwd()
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	for path, op := range opsByPath {
		opRef := filepath.Join(dirHandle.Ref(), path)
		if filepath.IsAbs(opRef) {
			// make absolute paths relative
			relOpRef, err := filepath.Rel(cwd, opRef)
			if nil != err {
				ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
				return // support fake exiter
			}

			opRef = strings.TrimPrefix(relOpRef, ".opspec/")
		}

		fmt.Fprintf(_tabWriter, "%v\t%v", opRef, op.Description)
		fmt.Fprintln(_tabWriter)

	}
}
