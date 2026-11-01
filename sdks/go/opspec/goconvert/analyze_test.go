package goconvert

import (
	"context"
	"github.com/opctl/opctl/sdks/go/model"
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
