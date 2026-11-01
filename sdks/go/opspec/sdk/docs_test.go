package sdk

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

var (
	goFence     = regexp.MustCompile("(?ms)^```go([^\n]*)\n(.*?)^```$")
	fenceTitle  = regexp.MustCompile(`title="([^"]+)"`)
	packageLine = regexp.MustCompile(`(?m)^package \w+$`)
)

// TestDocsExamplesCompile builds every Go code block in the website docs
// against this SDK. A block with a package clause is a whole file. Other
// blocks hold declarations and are wrapped in package main. Blocks titled
// with the same directory form one package, in a module named
// example.com/myapp.
func TestDocsExamplesCompile(t *testing.T) {
	if _, err := exec.LookPath("go"); err != nil {
		t.Skip("go is not on PATH")
	}
	sdk, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	docs := filepath.Join(sdk, "..", "..", "..", "..", "website", "docs")
	module := t.TempDir()
	goMod := "module example.com/myapp\n\ngo 1.21\n\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\n\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => " + sdk + "\n"
	if err := os.WriteFile(filepath.Join(module, "go.mod"), []byte(goMod), 0600); err != nil {
		t.Fatal(err)
	}
	// These interfaces are rendered and checked against shared describe fixtures
	// by goop's golden tests, without adding a dependency to this SDK module.
	for _, name := range []string{"greet", "merge"} {
		data, err := os.ReadFile(filepath.Join(sdk, "internal", "stubgolden", "doc"+name, "stub.go"))
		if err != nil {
			t.Fatal(err)
		}
		data = bytes.Replace(data, []byte("package doc"+name), []byte("package "+name), 1)
		dir := filepath.Join(module, "opctlgen", name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "stub.go"), data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	files := map[string]string{}
	dirs := map[string]bool{}
	err = filepath.WalkDir(docs, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".md" {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(docs, path)
		for _, match := range goFence.FindAllSubmatchIndex(data, -1) {
			meta, block := string(data[match[2]:match[3]]), string(data[match[4]:match[5]])
			line := strings.Count(string(data[:match[4]]), "\n") + 1
			name := fmt.Sprintf("example%d/op.go", len(files)+1)
			if title := fenceTitle.FindStringSubmatch(meta); title != nil {
				name = title[1]
			}
			if previous, ok := files[name]; ok {
				return fmt.Errorf("%s:%d: %s is also defined at %s", rel, line, name, previous)
			}
			files[name] = fmt.Sprintf("%s:%d", rel, line)
			source := fmt.Sprintf("//line %s:%d\n%s", rel, line, block)
			if !packageLine.MatchString(block) {
				source = "package main\n\nimport opctl \"github.com/opctl/opctl/sdks/go/opspec/sdk\"\n\nvar _ = opctl.Container\n\n" + source
				if !strings.Contains(block, "func main()") {
					source += "\nfunc main() {}\n"
				}
			}
			target := filepath.Join(module, filepath.FromSlash(name))
			if err := os.MkdirAll(filepath.Dir(target), 0700); err != nil {
				return err
			}
			if err := os.WriteFile(target, []byte(source), 0600); err != nil {
				return err
			}
			dirs["./"+filepath.ToSlash(filepath.Dir(name))] = true
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) < 10 {
		t.Fatalf("found %d Go code blocks in %s; expected the Go op guide and reference", len(files), docs)
	}
	packages := make([]string, 0, len(dirs))
	for dir := range dirs {
		packages = append(packages, dir)
	}
	sort.Strings(packages)
	cmd := exec.Command("go", append([]string{"vet"}, packages...)...)
	cmd.Dir = module
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly", "GOWORK=off", "GOTOOLCHAIN=local", "CGO_ENABLED=0")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("docs examples don't compile: %v\n%s", err, output)
	}
}
