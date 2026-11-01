package sdk

import (
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// Op overrides the op's name, which defaults to its directory name, and its
// description, which defaults to the package doc comment. Empty fields keep the default.
type Op struct {
	Name        string
	Description string
}

func overrides(op []Op) (Op, error) {
	if len(op) > 1 {
		return Op{}, errors.New("pass at most one opctl.Op")
	}
	if len(op) == 0 {
		return Op{}, nil
	}
	return op[0], nil
}

// packageDoc renders the doc comment of the main package in the working
// directory as markdown. go/build applies build constraints and skips tests.
func packageDoc() (string, error) {
	pkg, err := build.ImportDir(".", 0)
	var noGo *build.NoGoError
	if errors.As(err, &noGo) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("read package doc comment: %w", err)
	}
	if pkg.Name != "main" {
		return "", nil
	}
	fset := token.NewFileSet()
	files := make([]*ast.File, 0, len(pkg.GoFiles))
	for _, name := range pkg.GoFiles {
		file, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, name), nil, parser.PackageClauseOnly|parser.ParseComments)
		if err != nil {
			return "", fmt.Errorf("read package doc comment: %w", err)
		}
		files = append(files, file)
	}
	docs, err := doc.NewFromFiles(fset, files, "main")
	if err != nil {
		return "", fmt.Errorf("read package doc comment: %w", err)
	}
	return strings.TrimRight(string(docs.Markdown(docs.Doc)), "\n"), nil
}
