package core

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// Lser exposes the "ls" command
type Lser interface {
	Ls(
		ctx context.Context,
		dirRef string,
	) error
}

// newLser returns an initialized "ls" command
func newLser(
	cliOutput clioutput.CliOutput,
	dataResolver dataresolver.DataResolver,
) Lser {
	return _lsInvoker{
		cliOutput:    cliOutput,
		dataResolver: dataResolver,
		writer:       os.Stdout,
	}
}

type _lsInvoker struct {
	cliOutput    clioutput.CliOutput
	dataResolver dataresolver.DataResolver
	writer       io.Writer
}

func (ivkr _lsInvoker) Ls(
	ctx context.Context,
	dirRef string,
) error {
	_tabWriter := new(tabwriter.Writer)
	defer _tabWriter.Flush()
	_tabWriter.Init(ivkr.writer, 0, 8, 1, '\t', 0)

	fmt.Fprintln(_tabWriter, "REF\tDESCRIPTION")

	dirHandle, err := ivkr.dataResolver.Resolve(
		dirRef,
		nil,
	)
	if nil != err {
		return err
	}

	opsByPath, err := opspec.List(
		ctx,
		dirHandle,
	)
	if nil != err {
		return err
	}

	cwd, err := os.Getwd()
	if nil != err {
		return err
	}

	for path, op := range opsByPath {
		opRef := filepath.Join(dirHandle.Ref(), path)
		if filepath.IsAbs(opRef) {
			// make absolute paths relative
			relOpRef, err := filepath.Rel(cwd, opRef)
			if nil != err {
				return err
			}

			opRef = strings.TrimPrefix(relOpRef, ".opspec/")
		}

		fmt.Fprintf(_tabWriter, "%v\t%v", opRef, op.Description)
		fmt.Fprintln(_tabWriter)
	}

	return nil
}
