package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/tabwriter"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/opspath"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// ls implements "ls" command
func ls(
	ctx context.Context,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeConfig local.NodeConfig,
	dirRef string,
) error {
	np, err := local.New(nodeConfig)
	if err != nil {
		return err
	}

	node, err := np.CreateNodeIfNotExists(ctx)
	if err != nil {
		return err
	}

	opDirRefs, err := opspath.Get(
		ctx,
		dirRef,
		node,
	)
	if err != nil {
		return err
	}

	dataResolver := dataresolver.New(
		cliParamSatisfier,
		node,
	)

	_tabWriter := new(tabwriter.Writer)
	defer _tabWriter.Flush()
	_tabWriter.Init(os.Stdout, 0, 8, 1, '\t', 0)

	fmt.Fprintln(_tabWriter, "REF\tDESCRIPTION")

	for _, dirRef := range opDirRefs {

		dirHandle, err := dataResolver.Resolve(
			ctx,
			dirRef,
			nil,
		)
		if err != nil {
			return err
		}

		opsByPath, err := opspec.List(
			ctx,
			dirHandle,
		)
		if err != nil {
			return err
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		opdirRegexp := regexp.MustCompile(
			fmt.Sprintf(
				"^((%s)|(%s))?/",
				filepath.Join(
					model.DotOpctlDirName,
					//model.OpsDirName,
				),
				// gracefully handle deprecated ops path
				model.DotOpspecDirName,
			),
		)
		remoteRegexp := regexp.MustCompile(
			fmt.Sprintf(
				`^.*#\d\.\d\.\d/((%s)|(%s))?/`,
				filepath.Join(
					model.DotOpctlDirName,
					//model.OpsDirName,
				),
				// gracefully handle deprecated ops path
				model.DotOpspecDirName,
			),
		)

		for path, op := range opsByPath {
			opRef := filepath.Join(dirHandle.Ref(), path)
			if filepath.IsAbs(opRef) {
				// make absolute paths relative
				relOpRef, err := filepath.Rel(
					cwd,
					opRef,
				)
				if err != nil {
					return err
				}

				if parts := strings.Split(relOpRef, "#"); len(parts) == 2 {
					opRef = parts[1]
				}

				opRef = opdirRegexp.ReplaceAllString(relOpRef, "")
			} else {
				opRef = remoteRegexp.ReplaceAllString(opRef, "")
			}

			description := op.Description
			if strings.TrimSpace(description) == "" {
				description = "-"
			}

			scanner := bufio.NewScanner(strings.NewReader(description))
			if scanner.Scan() {
				// first line of description, add the op ref
				fmt.Fprintf(_tabWriter, "%v\t%v", opRef, scanner.Text())
			}
			for scanner.Scan() {
				// subsequent lines, don't add the op ref but let the description span multiple lines
				fmt.Fprintf(_tabWriter, "\n\t%v", scanner.Text())
			}
			fmt.Fprintln(_tabWriter)
		}
	}

	return nil
}
