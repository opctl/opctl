package goconvert

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRenderGraph(t *testing.T) {
	for _, fixture := range []string{"streams", "typed-outputs/pass"} {
		t.Run(fixture, func(t *testing.T) {
			dir, err := filepath.Abs("../../../../test-suite/run/container/object/" + fixture)
			if err != nil {
				t.Fatal(err)
			}
			source, err := os.ReadFile(filepath.Join(dir, "op.yml"))
			if err != nil {
				t.Fatal(err)
			}
			a, err := Analyze(context.Background(), dir, source, nil)
			if err != nil || len(a.Refusals) > 0 {
				t.Fatalf("%v %+v", err, a)
			}
			rendered, err := Render(a, RenderOptions{})
			if err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join("testdata", filepath.Base(fixture)+".go.golden")
			if os.Getenv("UPDATE_GOLDEN") == "1" {
				if err = os.MkdirAll("testdata", 0755); err != nil {
					t.Fatal(err)
				}
				if err = os.WriteFile(golden, rendered.Source, 0600); err != nil {
					t.Fatal(err)
				}
			}
			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(want, rendered.Source) {
				t.Fatalf("source differs from %s:\n%s", golden, rendered.Source)
			}
			planned, err := goop.ParsePlan(runCode(t, rendered, "plan", nil))
			if err != nil {
				t.Fatal(err)
			}
			if planned.Run.Serial == nil {
				t.Fatalf("not a serial graph: %+v", planned.Run)
			}
			if fixture == "typed-outputs/pass" {
				if len(planned.OutputBindings) != 8 {
					t.Fatalf("lost outputs: %+v", planned.OutputBindings)
				}
				first := (*planned.Run.Serial)[0].Container
				if len(first.Outputs) != 6 || len(first.OutputParams) != 6 || first.Stdout == "" || first.Stderr == "" {
					t.Fatalf("lost producers: %+v", first)
				}
			}
		})
	}
}
func runCode(t *testing.T, rendered *Rendered, phase string, inputs map[string]any) []byte {
	t.Helper()
	dir := t.TempDir()
	sdk, err := filepath.Abs("../sdk")
	if err != nil {
		t.Fatal(err)
	}
	files := map[string][]byte{"op.go": rendered.Source, "go.mod": []byte("module ops.local/test\ngo 1.25\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => " + sdk + "\n")}
	for name, data := range rendered.Stubs {
		files[filepath.Join("opctlgen", name, "stub.go")] = data
	}
	request, _ := json.Marshal(goop.NewRequest(phase, inputs))
	files["request.json"] = request
	for name, data := range files {
		path := filepath.Join(dir, name)
		if err = os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err = os.WriteFile(path, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "run", ".", phase)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly", "GOWORK=off", "GOTOOLCHAIN=local", "OPCTL_BOOTSTRAP_REQUEST="+filepath.Join(dir, "request.json"))
	out, err := cmd.Output()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			t.Fatalf("%v %s\n%s", err, e.Stderr, rendered.Source)
		}
		t.Fatal(err)
	}
	return out
}
