package goconvert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func interfaceJSON(spec *model.OpSpec) string {
	raw, _ := json.Marshal(spec)
	var doc map[string]any
	_ = json.Unmarshal(raw, &doc)
	delete(doc, "run")
	delete(doc, "version")
	for _, group := range []string{"inputs", "outputs"} {
		params, _ := doc[group].(map[string]any)
		for _, v := range params {
			p := v.(map[string]any)
			for _, kind := range []string{"array", "object", "string", "number", "boolean", "dir", "file", "socket"} {
				if body, ok := p[kind].(map[string]any); ok {
					if p["description"] == nil && body["description"] != nil {
						p["description"] = body["description"]
					}
					delete(body, "description")
					if c, ok := body["constraints"].(map[string]any); ok && len(c) == 0 {
						delete(body, "constraints")
					}
					if body["isSecret"] == false {
						delete(body, "isSecret")
					}
				}
			}
		}
	}
	out, _ := json.Marshal(doc)
	return string(out)
}
func TestInterfaceCorpus(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles converted interfaces across the corpus")
	}
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}
	sdk := filepath.Join(root, "sdks/go/opspec/sdk")
	module := t.TempDir()
	bin := filepath.Join(module, "bin")
	if err = os.Mkdir(bin, 0755); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(module, "go.mod"), []byte("module ops.local/corpus\ngo 1.25\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => "+sdk+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	type testCase struct {
		dir, source string
		spec        *model.OpSpec
	}
	cases := []testCase{}
	describe := func(context.Context, string) (*model.OpSpec, error) { return &model.OpSpec{}, nil }
	for _, group := range []string{"test-suite", "examples", ".opspec"} {
		err = filepath.WalkDir(filepath.Join(root, group), func(path string, entry os.DirEntry, err error) error {
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
			if err != nil || len(a.Refusals) > 0 {
				return nil
			}
			name := fmt.Sprintf("case%04d", len(cases)+1)
			dir := filepath.Join(module, name)
			if err = os.Mkdir(dir, 0755); err != nil {
				return err
			}
			a.Dir = dir
			code, err := RenderInterface(a)
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			if err = os.WriteFile(filepath.Join(dir, "op.go"), code, 0600); err != nil {
				return err
			}
			cases = append(cases, testCase{dir, path, a.Spec})
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(cases) < 20 {
		t.Fatalf("only %d interfaces", len(cases))
	}
	cmd := exec.Command("go", "build", "-o", bin+string(os.PathSeparator), "./...")
	cmd.Dir = module
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly -buildvcs=false", "GOWORK=off", "GOTOOLCHAIN=local")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build interfaces: %v\n%s", err, out)
	}
	request, _ := json.Marshal(goop.NewRequest("describe", nil))
	requestPath := filepath.Join(module, "request.json")
	if err = os.WriteFile(requestPath, request, 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range cases {
		cmd := exec.Command(filepath.Join(bin, filepath.Base(tc.dir)), "describe")
		cmd.Dir = tc.dir
		cmd.Env = append(os.Environ(), "OPCTL_BOOTSTRAP_REQUEST="+requestPath)
		out, err := cmd.Output()
		if err != nil {
			t.Fatalf("%s: %v", tc.source, err)
		}
		var response struct{ Opspec *model.OpSpec }
		if err = json.Unmarshal(out, &response); err != nil {
			t.Fatal(err)
		}
		if got, want := interfaceJSON(response.Opspec), interfaceJSON(tc.spec); got != want {
			t.Fatalf("%s interface mismatch:\n%s\n%s", tc.source, got, want)
		}
	}
	t.Logf("compiled and described %d converted interfaces", len(cases))
}
