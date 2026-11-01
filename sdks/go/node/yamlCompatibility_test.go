package node

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

func TestFailedYAMLOpRetainsBoundOutputs(t *testing.T) {
	dir := t.TempDir()
	declaration := "name: invalid-output\noutputs:\n  result:\n    string:\n      constraints: {minLength: 3}\nrun:\n  container:\n    image: {ref: alpine}\n"
	if err := os.WriteFile(filepath.Join(dir, "op.yml"), []byte(declaration), 0600); err != nil {
		t.Fatal(err)
	}
	child := new(fakes.FakeCaller)
	text := "x"
	child.CallReturns(map[string]*model.Value{"result": {String: &text}}, nil)
	bus := new(fakes.FakePubSub)
	events := make(chan model.Event)
	close(events)
	bus.SubscribeReturns(events, nil)
	caller := _caller{pubSub: bus, opCaller: _opCaller{caller: child, dataDirPath: dir}, dataDirPath: dir}
	result, err := caller.Call(context.Background(), "root", nil, &model.CallSpec{Op: &model.OpCallSpec{Ref: dir, PullCreds: &model.CredsSpec{}, Outputs: map[string]string{"result": "$(renamed)"}}}, dir, nil, "root")
	if err == nil {
		t.Fatal("invalid output passed")
	}
	if result["renamed"] == nil || result["renamed"].String == nil || *result["renamed"].String != text {
		t.Fatal("failed YAML op dropped its bound output")
	}
	end := bus.PublishArgsForCall(bus.PublishCallCount() - 1).CallEnded
	if end == nil || end.Outputs["renamed"] == nil || *end.Outputs["renamed"].String != text {
		t.Fatal("CallEnded dropped the legacy bound output")
	}
}

func TestFailedYAMLOutputsAreScopedToSubtree(t *testing.T) {
	for _, sibling := range []string{"none", "results", "go"} {
		for _, descendant := range []string{"legacy", "results", "go"} {
			t.Run(sibling+"/"+descendant, func(t *testing.T) {
				dir := t.TempDir()
				write := func(name, file, declaration string) string {
					t.Helper()
					path := filepath.Join(dir, name)
					if err := os.MkdirAll(path, 0700); err != nil {
						t.Fatal(err)
					}
					if err := os.WriteFile(filepath.Join(path, file), []byte(declaration), 0600); err != nil {
						t.Fatal(err)
					}
					return path
				}
				leaf := func(name, kind string) string {
					t.Helper()
					if kind == "go" {
						path := write(name, "op.go", "package main\n")
						write(name, "go.mod", "module example.com/leaf\ngo 1.21\n")
						return path
					}
					declaration := "name: leaf\noutputs:\n  result:\n    string: {}\nrun:\n  container:\n    image: {ref: alpine}\n"
					if kind == "results" {
						declaration += "    stdout: $(result)\n"
					}
					return write(name, "op.yml", declaration)
				}
				leaf("leaf", descendant)
				write("wrapper", "op.yml", "name: wrapper\noutputs:\n  result:\n    string: {}\nrun:\n  op:\n    ref: $(../leaf)\n    outputs: {result: $(result)}\n")
				parent := write("parent", "op.yml", "name: invalid-output\noutputs:\n  result:\n    string:\n      constraints: {minLength: 3}\nrun:\n  op:\n    ref: $(../wrapper)\n    outputs: {result: $(result)}\n")
				bus := new(fakes.FakePubSub)
				runtime := new(fakes.FakeContainerCaller)
				text := "x"
				runtime.CallReturns(map[string]*model.Value{"result": {String: &text}}, nil)
				caller := newCaller(runtime, dir, bus, &bootstrapSnapshotStore{}, &goBootstrapCache{}).(*_caller)
				bootstrap := new(fakes.FakeBootstrapper)
				desc := &model.OpSpec{Name: "leaf", Outputs: map[string]*model.ParamSpec{"result": {String: &model.StringParamSpec{}}}}
				bootstrap.DescribeReturns(desc, nil)
				bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}}}}, OutputBindings: map[string]string{"result": "$(result)"}}, nil)
				caller.opCaller = _opCaller{caller: caller, dataDirPath: dir, bootstrapper: bootstrap}
				ctx := withPrivateCompletions(context.Background())
				defer clearPrivateCompletions(ctx)
				invoke := func(id, path string) (map[string]*model.Value, error) {
					return caller.Call(ctx, id, nil, &model.CallSpec{Op: &model.OpCallSpec{Ref: path, PullCreds: &model.CredsSpec{}, Outputs: map[string]string{"result": "$(renamed)"}}}, dir, nil, "root")
				}
				if sibling != "none" {
					if _, err := invoke("sibling", leaf("sibling", sibling)); err != nil {
						t.Fatal(err)
					}
					if !privateResultsActive(ctx) {
						t.Fatal("sibling did not activate private results")
					}
				}
				result, err := invoke("parent", parent)
				if err == nil || !strings.Contains(err.Error(), "result: String length must be greater than or equal to 3") {
					t.Fatalf("expected parent output validation failure, got %v", err)
				}
				end := bus.PublishArgsForCall(bus.PublishCallCount() - 1).CallEnded
				if descendant == "legacy" {
					if result["renamed"] == nil || result["renamed"].String == nil || *result["renamed"].String != text {
						t.Fatal("unrelated sibling dropped legacy failure outputs")
					}
					if end == nil || end.Outputs["renamed"] == nil {
						t.Fatal("legacy failure event lost its output")
					}
				} else if len(result) != 0 || end == nil || len(end.Outputs) != 0 {
					t.Fatal("a subtree with new results returned partial failure outputs")
				}
			})
		}
	}
}
