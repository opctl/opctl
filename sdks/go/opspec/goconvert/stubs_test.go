package goconvert

import (
	"context"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConvertedTypedStub(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "dev")
	if err := os.Mkdir(dir, 0755); err != nil {
		t.Fatal(err)
	}
	source, err := os.ReadFile("testdata/go-service.yml")
	if err != nil {
		t.Fatal(err)
	}
	calleeSource, err := os.ReadFile("../../../../examples/run-a-go-service/.opspec/mysql/op.yml")
	if err != nil {
		t.Fatal(err)
	}
	callee, err := opfile.Unmarshal(calleeSource)
	if err != nil {
		t.Fatal(err)
	}
	a, err := Analyze(context.Background(), dir, source, func(context.Context, string) (*model.OpSpec, error) { return callee, nil })
	if err != nil || len(a.Refusals) > 0 {
		t.Fatalf("%v %+v", err, a.Refusals)
	}
	rendered, err := Render(a, RenderOptions{ModuleRoot: dir, ModulePath: "ops.local/test"})
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"mysql.Call(mysql.Inputs{", "//opctl:op ../mysql", "MYSQLUSER:", "opctl.RawRef[opctl.Dir]", "TODO(opctl op convert)"} {
		if !strings.Contains(string(rendered.Source), want) {
			t.Fatalf("missing %q:\n%s", want, rendered.Source)
		}
	}
	golden := "testdata/go-service.go.golden"
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
		t.Fatalf("service golden changed:\n%s", rendered.Source)
	}
	planned, err := goop.ParsePlan(runCode(t, rendered, "plan", map[string]any{"goBuildCachePath": "Library/Caches/go-build", "MYSQL_USER": "testuser", "MYSQL_PASSWORD": "testpassword", "MYSQL_DATABASE": "testapp", "MYSQL_HOST": "mysql"}))
	if err != nil {
		t.Fatal(err)
	}
	described, err := goop.ParseDescribe(runCode(t, rendered, "describe", nil, "dev"))
	if err != nil {
		t.Fatal(err)
	}
	if err = Compare(a, described, planned, map[string]any{"goBuildCachePath": "Library/Caches/go-build", "MYSQL_USER": "testuser", "MYSQL_PASSWORD": "testpassword", "MYSQL_DATABASE": "testapp", "MYSQL_HOST": "mysql"}); err != nil {
		t.Fatal(err)
	}
	if planned.Run.Parallel == nil || (*planned.Run.Parallel)[0].Op.Ref != "../mysql" {
		t.Fatalf("wrong stub ref: %+v", planned.Run)
	}
}
func TestStubFallbackOnArrayType(t *testing.T) {
	dir := t.TempDir()
	source := []byte("name: fallback\ninputs: {tags: {array: {default: [x]}}}\nrun: {op: {ref: ../child, inputs: {tags: '$(tags)'}}}\n")
	callee := &model.OpSpec{Inputs: map[string]*model.ParamSpec{"tags": {Array: &model.ArrayParamSpec{Constraints: model.Constraints{"items": map[string]any{"type": "string"}}}}}}
	a, err := Analyze(context.Background(), dir, source, func(context.Context, string) (*model.OpSpec, error) { return callee, nil })
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := Render(a, RenderOptions{ModuleRoot: dir, ModulePath: "ops.local/test"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(rendered.Source), "CallOp, not a stub:") || !strings.Contains(string(rendered.Source), "Array[any]") || len(rendered.Stubs) != 0 {
		t.Fatal(string(rendered.Source))
	}
	_ = runCode(t, rendered, "plan", map[string]any{"tags": []any{"x"}})
}

func TestConvertedStubNames(t *testing.T) {
	dir := t.TempDir()
	source := []byte("name: names\nrun:\n  serial:\n  - op: {ref: ../one/build}\n  - op: {ref: ../two/build}\n  - op: {ref: ../plan}\n  - op: {ref: ../one/build}\n")
	a, err := Analyze(context.Background(), dir, source, func(context.Context, string) (*model.OpSpec, error) { return &model.OpSpec{}, nil })
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := Render(a, RenderOptions{ModuleRoot: dir, ModulePath: "ops.local/test"})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"onebuild", "twobuild", "plan2"} {
		if rendered.Stubs[name] == nil {
			t.Fatalf("missing %s in %v", name, rendered.Stubs)
		}
	}
	if len(rendered.Stubs) != 3 {
		t.Fatal(rendered.Stubs)
	}
	_ = runCode(t, rendered, "plan", nil)
}
