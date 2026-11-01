package goconvert

import (
	"context"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderControl(t *testing.T) {
	source := []byte("name: control\ninputs:\n  env: {string: {default: dev}}\n  count: {number: {default: 2}}\n  regions: {array: {default: [east, west], constraints: {items: {type: string}}}}\noutputs: {echo: {string: {default: none}}}\nrun:\n  serial:\n    - container:\n        image: {ref: alpine}\n        stdout: $(echo)\n    - if: [{eq: ['$(env)', prod]}, {gt: ['$(count)', 1]}]\n      container:\n        image: {ref: alpine}\n        stdout: $(echo)\n    - serialLoop:\n        range: $(regions)\n        vars: {index: '$(i)', value: '$(region)'}\n        run:\n          container:\n            image: {ref: alpine}\n            dirs: {/scratch: '$(scratch)'}\n            envVars: {REGION: '$(region)', INDEX: '$(i)'}\n    - parallelLoop:\n        range: [a, b]\n        vars: {value: '$(name)'}\n        run:\n          container:\n            image: {ref: alpine}\n            dirs: {/scratch: '$(fresh)'}\n            envVars: {NAME: '$(name)'}\n")
	a, err := Analyze(context.Background(), t.TempDir(), source, nil)
	if err != nil || len(a.Refusals) > 0 {
		t.Fatalf("%v %+v", err, a.Refusals)
	}
	rendered, err := Render(a, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"if ", "for ", "in.Regions.Value()", "opctl.NewDir()"} {
		if !strings.Contains(string(rendered.Source), want) {
			t.Fatalf("missing %s:\n%s", want, rendered.Source)
		}
	}
	golden := "testdata/control.go.golden"
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err = os.WriteFile(golden, rendered.Source, 0644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if string(want) != string(rendered.Source) {
		t.Fatalf("control golden changed:\n%s", rendered.Source)
	}
	for _, env := range []string{"dev", "prod"} {
		plan, err := goop.ParsePlan(runCode(t, rendered, "plan", map[string]any{"env": env, "count": 2, "regions": []any{"east", "west"}}))
		if err != nil {
			t.Fatal(err)
		}
		first := (*plan.Run.Serial)[0].Container
		if env == "dev" && plan.OutputBindings["echo"] != first.Stdout {
			t.Fatalf("false branch lost the preceding output: %+v", plan.OutputBindings)
		}
	}
}
func TestRenderObjectLoop(t *testing.T) {
	source := []byte("name: object-loop\nrun:\n  serialLoop:\n    range: {b: west, a: east}\n    vars: {index: '$(i)', key: '$(key)', value: '$(value)'}\n    run:\n      container:\n        image: {ref: alpine}\n        cmd: ['echo', '$(i):$(key):$(value)']\n")
	dir := filepath.Join(t.TempDir(), "object-loop")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	a, err := Analyze(context.Background(), dir, source, nil)
	if err != nil || len(a.Refusals) > 0 {
		t.Fatalf("%v %+v", err, a.Refusals)
	}
	rendered, err := Render(a, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := goop.ParsePlan(runCode(t, rendered, "plan", nil))
	if err != nil {
		t.Fatal(err)
	}
	first := (*plan.Run.Serial)[0].Container.Cmd.([]any)
	if first[1] != "0:a:east" {
		t.Fatalf("object keys weren't sorted: %+v", first)
	}
}
