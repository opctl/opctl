package goconvert

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/internal/gooptest"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanCorpus(t *testing.T) {
	if testing.Short() {
		t.Skip("builds and verifies converted graphs across the op corpus")
	}
	root, err := filepath.Abs("../../../..")
	if err != nil {
		t.Fatal(err)
	}
	module := t.TempDir()
	bin := filepath.Join(module, "bin")
	if err = os.Mkdir(bin, 0755); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(module, "go.mod"), []byte("module ops.local/corpus\ngo 1.25\nrequire github.com/opctl/opctl/sdks/go/opspec/sdk v0.0.0\nreplace github.com/opctl/opctl/sdks/go/opspec/sdk => "+filepath.Join(root, "sdks/go/opspec/sdk")+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	describe := func(ctx context.Context, ref string) (*model.OpSpec, error) {
		if strings.Contains(ref, "#") {
			return nil, fmt.Errorf("remote callee needs node describe: %s", ref)
		}
		data, err := os.ReadFile(filepath.Join(ref, "op.yml"))
		if err == nil {
			return opfile.Unmarshal(data)
		}
		if _, e := os.Stat(filepath.Join(ref, "op.go")); e == nil {
			out, e := gooptest.Phase(ref, "describe", nil)
			if e != nil {
				return nil, e
			}
			return goop.ParseDescribe(out)
		}
		return nil, err
	}
	type testCase struct {
		a               *Analysis
		name, path, dir string
		sets            []map[string]any
	}
	cases := []testCase{}
	refused := 0
	reasons := map[string]int{}
	for _, group := range []string{"test-suite", "examples", ".opspec", "cli/.opspec"} {
		err = filepath.WalkDir(filepath.Join(root, group), func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || entry.Name() != "op.yml" {
				return nil
			}
			source, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			a, err := Analyze(context.Background(), filepath.Dir(path), source, describe)
			if err != nil {
				refused++
				reasons["parse or describe unavailable"]++
				t.Logf("REFUSED %s: %v", strings.TrimPrefix(path, root+"/"), err)
				return nil
			}
			if len(a.Refusals) > 0 {
				refused++
				for _, r := range a.Refusals {
					reasons[r.Reason]++
					t.Logf("REFUSED %s %s", strings.TrimPrefix(path, root+"/"), r)
				}
				return nil
			}
			sets, err := InputSets(a.Spec, a.Dir, nil)
			if err != nil {
				refused++
				reasons[err.Error()]++
				t.Logf("REFUSED %s: %v", path, err)
				return nil
			}
			name := fmt.Sprintf("case%04d", len(cases)+1)
			dir := filepath.Join(module, name)
			if err = os.Mkdir(dir, 0755); err != nil {
				return err
			}
			staged := *a
			staged.Dir = dir // Explicit metadata preserves the original name in the temporary package.
			rendered, err := Render(&staged, RenderOptions{})
			if err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
			if err = os.WriteFile(filepath.Join(dir, "op.go"), rendered.Source, 0600); err != nil {
				return err
			}
			cases = append(cases, testCase{a, name, strings.TrimPrefix(path, root+"/"), dir, sets})
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command("go", "build", "-o", bin+string(os.PathSeparator), "./...")
	cmd.Dir = module
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=readonly -buildvcs=false", "GOWORK=off", "GOTOOLCHAIN=local")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build corpus: %v\n%s", err, out)
	}
	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			phase := func(name string, inputs map[string]any) []byte {
				request, _ := json.Marshal(goop.NewRequest(name, inputs))
				requestPath := filepath.Join(module, "request.json")
				if err := os.WriteFile(requestPath, request, 0600); err != nil {
					t.Fatal(err)
				}
				cmd := exec.Command(filepath.Join(bin, tc.name), name)
				cmd.Dir = tc.dir
				cmd.Env = append(os.Environ(), "OPCTL_BOOTSTRAP_REQUEST="+requestPath)
				out, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("%s: %v %s", name, err, out)
				}
				return out
			}
			desc, err := goop.ParseDescribe(phase("describe", nil))
			if err != nil {
				t.Fatal(err)
			}
			for _, inputs := range tc.sets {
				plan, err := goop.ParsePlan(phase("plan", inputs))
				if err != nil {
					t.Fatal(err)
				}
				if err = Compare(tc.a, desc, plan, inputs); err != nil {
					t.Fatalf("inputs %v: %v", inputs, err)
				}
			}
		})
	}
	t.Logf("CORPUS converted=%d refused=%d", len(cases), refused)
	for _, reason := range keys(reasons) {
		t.Logf("REFUSAL %d %s", reasons[reason], reason)
	}
	if len(cases) < 50 {
		t.Fatalf("only %d converted graphs", len(cases))
	}
}
