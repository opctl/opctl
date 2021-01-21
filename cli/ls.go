package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/opspec"
)

func ls(node node.Node, dirRef string, dataResolver dataresolver.DataResolver) error {
	_tabWriter := new(tabwriter.Writer)
	defer _tabWriter.Flush()
	_tabWriter.Init(os.Stdout, 0, 8, 1, '\t', 0)

	fmt.Fprintln(_tabWriter, "REF\tDESCRIPTION")

	dirHandle, err := dataResolver.Resolve(
		dirRef,
		nil,
	)
	if nil != err {
		return err
	}

	opsByPath, err := opspec.List(
		context.TODO(),
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
