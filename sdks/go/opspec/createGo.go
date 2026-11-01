package opspec

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"golang.org/x/mod/module"
)

const goSDKModule = "github.com/opctl/opctl/sdks/go/opspec/sdk"

// GoSDK identifies a published Go op SDK release and its go.sum hashes.
type GoSDK struct {
	Version  string
	Sum      string
	GoModSum string
}

// Published reports whether scaffolds can require this release without a replace directive.
func (sdk GoSDK) Published() bool {
	return sdk.Version != "" && sdk.Sum != "" && sdk.GoModSum != ""
}

const goOpTemplate = `package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Who opctl.String ` + "`" + `opctl:"who,description='Who to greet',default=world"` + "`" + `
}

func plan(in Inputs) opctl.Call {
	return opctl.Container("alpine:3").
		Cmd("sh", "-ce", ` + "`" + `echo "hello $WHO"` + "`" + `).
		EnvFrom("WHO", in.Who)
}

func main() { opctl.Main(plan) }
`

// CreateGo creates a Go op in its own module at path. It never overwrites an
// existing op or module file.
func CreateGo(path, name, description string, sdk GoSDK) error {
	if err := module.CheckImportPath(name); err != nil {
		return fmt.Errorf("op name %q is not a valid Go module path", name)
	}

	opGo := goOpTemplate
	if description != "" {
		opGo = "// " + strings.ReplaceAll(description, "\n", "\n// ") + "\n" + opGo
		opGo = strings.ReplaceAll(opGo, "// \n", "//\n")
	}
	version := "v0.0.0"
	if sdk.Published() {
		version = sdk.Version
	}
	files := []struct{ name, content string }{
		{"op.go", opGo},
		{"go.mod", fmt.Sprintf("module %s\n\ngo %s\n\nrequire %s %s\n", name, goDirective(), goSDKModule, version)},
	}
	if sdk.Published() {
		files = append(files, struct{ name, content string }{"go.sum", fmt.Sprintf("%s %s %s\n%s %s/go.mod %s\n", goSDKModule, version, sdk.Sum, goSDKModule, version, sdk.GoModSum)})
	}

	for _, existing := range []string{opfile.FileName, "op.go", "go.mod", "go.sum"} {
		if _, err := os.Lstat(filepath.Join(path, existing)); err == nil {
			return fmt.Errorf("%s already exists in %s", existing, path)
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	if err := os.MkdirAll(path, 0777); err != nil {
		return err
	}
	for i, file := range files {
		if err := writeNew(filepath.Join(path, file.name), file.content); err != nil {
			for _, written := range files[:i] {
				os.Remove(filepath.Join(path, written.name))
			}
			return err
		}
	}
	return nil
}

func writeNew(path, content string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	if _, err := file.WriteString(content); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

// goDirective follows the toolchain that built opctl, which selects the bootstrap image.
func goDirective() string { return goDirectiveFor(runtime.Version()) }

var goRelease = regexp.MustCompile(`^go(\d+\.\d+)`)

func goDirectiveFor(version string) string {
	if match := goRelease.FindStringSubmatch(version); match != nil {
		return match[1]
	}
	return "1.25"
}
