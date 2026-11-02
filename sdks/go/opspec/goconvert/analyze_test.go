package goconvert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzeRefusals(t *testing.T) {
	for _, tc := range []struct{ name, source, want string }{
		{"ports", "name: example\nrun:\n  container:\n    image: {ref: alpine}\n    ports: {80: 80}\n", "ports"},
		{"initializer", "name: example\nrun:\n  container:\n    image: {ref: alpine}\n    files: {/x: {content: x}}\n", "initializers"},
		{"name", "name: example\nrun:\n  name: named\n  container:\n    image: {ref: alpine}\n", "call name"},
		{"unknown", "name: example\nrun:\n  container:\n    image: {ref: alpine}\n    envVars: {X: '$(missing)'}\n", "unknown reference"},
		{"result template", "name: example\nrun:\n  serial:\n    - container:\n        image: {ref: alpine}\n        stdout: $(text)\n    - container:\n        image: {ref: alpine}\n        cmd: ['echo', 'prefix $(text)']\n", "produced"},
		{"parallel sibling", "name: example\nrun:\n  parallel:\n    - container:\n        image: {ref: alpine}\n        stdout: $(text)\n    - container:\n        image: {ref: alpine}\n        envVars: {X: '$(text)'}\n", "unknown reference"},
		{"missing output", "name: example\noutputs: {result: {string: {}}}\nrun: {container: {image: {ref: alpine}}}\n", "unbound"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a, err := Analyze(context.Background(), t.TempDir(), []byte(tc.source), nil)
			if err != nil {
				t.Fatal(err)
			}
			found := false
			for _, r := range a.Refusals {
				if r.Line < 1 || r.Column < 1 {
					t.Fatalf("missing location: %+v", r)
				}
				if strings.Contains(r.Reason, tc.want) {
					found = true
				}
			}
			if !found {
				t.Fatalf("want %s: %+v", tc.want, a.Refusals)
			}
		})
	}
}
func TestAnalyzeSupported(t *testing.T) {
	source := []byte("name: example\ninputs:\n  token: {string: {isSecret: true}}\n  text: {string: {default: hi}}\noutputs: {result: {string: {}}}\nrun:\n  serial:\n    - container:\n        image: {ref: alpine}\n        dirs: {/scratch: '$(scratch)'}\n        envVars: {AUTH: 'Bearer $(token)', TEXT: 'prefix $(text)'}\n        stdout: $(result)\n    - container:\n        image: {ref: alpine}\n        envVars: {TEXT: '$(result)'}\n")
	a, err := Analyze(context.Background(), t.TempDir(), source, nil)
	if err != nil || len(a.Refusals) != 0 {
		t.Fatalf("%v %+v", err, a)
	}
}
func TestAnalysisCorpus(t *testing.T) {
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}
	total := 0
	describe := func(ctx context.Context, ref string) (*model.OpSpec, error) {
		return &model.OpSpec{Name: filepath.Base(ref)}, nil
	}
	for _, dir := range []string{"test-suite", "examples", ".opspec"} {
		err = filepath.WalkDir(filepath.Join(root, dir), func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || entry.Name() != "op.yml" {
				return nil
			}
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			a, err := Analyze(context.Background(), filepath.Dir(path), data, describe)
			if err != nil {
				return nil
			}
			total++
			for _, r := range a.Refusals {
				if r.Line < 1 || r.Column < 1 {
					t.Fatalf("%s: %+v", path, r)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if total < 250 {
		t.Fatalf("only %d parseable ops", total)
	}
}

func TestAnalyzeRejectsDirectiveInjection(t *testing.T) {
	for _, ref := range []string{"../child\n)\nfunc init() { println(1) }\nvar (_ = 1", "../child other", "../child\tother", "../child\rfoo", "../child\x00", "../child\u00a0other"} {
		t.Run(ref, func(t *testing.T) {
			source, err := json.Marshal(map[string]any{"name": "injection", "run": map[string]any{"op": map[string]any{"ref": ref}}})
			if err != nil {
				t.Fatal(err)
			}
			described := false
			a, err := Analyze(context.Background(), t.TempDir(), source, func(context.Context, string) (*model.OpSpec, error) { described = true; return &model.OpSpec{}, nil })
			if err != nil {
				t.Fatal(err)
			}
			if described {
				t.Fatal("unsafe ref reached describe")
			}
			if len(a.Refusals) == 0 || a.Refusals[0].Reason != "op refs can't contain whitespace or control characters" {
				t.Fatalf("refusals: %+v", a.Refusals)
			}
			if _, err = Render(a, RenderOptions{ModuleRoot: a.Dir, ModulePath: "ops.local/injection"}); err == nil {
				t.Fatal("render accepted unsafe ref")
			}
		})
	}
}

func TestRenderQuotesSourceText(t *testing.T) {
	payload := "\"`\n)\nfunc init() { println(1) }\nvar (_ = 1 //"
	source, err := json.Marshal(map[string]any{
		"name": payload, "description": payload,
		"inputs": map[string]any{"text": map[string]any{"description": payload, "string": map[string]any{"default": payload, "constraints": map[string]any{"enum": []any{payload}}}}},
		"run":    map[string]any{"description": payload, "container": map[string]any{"image": map[string]any{"ref": "alpine"}, "cmd": []any{payload}, "envVars": map[string]any{"TEXT": payload}, "files": map[string]any{"/file": "$(./file)"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	a, err := Analyze(context.Background(), t.TempDir(), source, nil)
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := Render(a, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	file, err := parser.ParseFile(token.NewFileSet(), "op.go", rendered.Source, 0)
	if err != nil {
		t.Fatal(err)
	}
	functions := 0
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			functions++
			if f.Name.Name != "main" && f.Name.Name != "plan" {
				t.Fatalf("injected function %s", f.Name.Name)
			}
		}
	}
	if functions != 2 {
		t.Fatalf("got %d functions", functions)
	}
}

func TestLoopRefusesOuterReadAndRebind(t *testing.T) {
	for _, loop := range []string{"serialLoop", "parallelLoop"} {
		for _, read := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/read=%t", loop, read), func(t *testing.T) {
				env := ""
				if read {
					env = ", envVars: {PREVIOUS: '$(text)'}"
				}
				source := fmt.Sprintf("name: loop-rebind\nrun:\n  serial:\n  - container: {image: {ref: alpine}, stdout: '$(text)'}\n  - %s:\n      range: [a, b]\n      run:\n        container: {image: {ref: alpine}%s, stdout: '$(text)'}\n", loop, env)
				a, err := Analyze(context.Background(), t.TempDir(), []byte(source), nil)
				if err != nil {
					t.Fatal(err)
				}
				refused := false
				for _, r := range a.Refusals {
					refused = refused || strings.Contains(r.Reason, "loop body reads and rebinds outer name text across iterations")
				}
				if refused != read {
					t.Fatalf("read=%t: %+v", read, a.Refusals)
				}
			})
		}
	}
}
