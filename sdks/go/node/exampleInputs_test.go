package node

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	callpkg "github.com/opctl/opctl/sdks/go/opspec/interpreter/call"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/inputs"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

func TestReactDevPassesHomeToGoService(t *testing.T) {
	ctx := context.Background()
	path, err := filepath.Abs("../../../examples/run-a-react-app/.opspec/dev")
	if err != nil {
		t.Fatal(err)
	}
	scratch := t.TempDir()
	dir, err := fs.New(scratch).TryResolve(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	spec, err := opfile.Get(ctx, dir)
	if err != nil {
		t.Fatal(err)
	}
	home := t.TempDir()
	scope, err := inputs.Interpret(map[string]interface{}{"HOME": "$(home)"}, spec.Inputs, path, map[string]*model.Value{"home": {Dir: &home}}, scratch)
	if err != nil {
		t.Fatal(err)
	}
	parent := filepath.Dir(path)
	scope["./"] = &model.Value{Dir: &path}
	scope["../"] = &model.Value{Dir: &parent}
	scope["/"] = &model.Value{Dir: &path}
	if spec.Run.Parallel == nil || len(*spec.Run.Parallel) != 2 {
		t.Fatal("expected the API and UI in parallel")
	}
	for _, child := range *spec.Run.Parallel {
		if child.Op == nil {
			continue
		}
		child.Op.PullCreds = &model.CredsSpec{}
		call, err := callpkg.Interpret(ctx, scope, child, "service", path, nil, "root", scratch)
		if err != nil {
			t.Fatal(err)
		}
		bootstrap := new(fakes.FakeBootstrapper)
		bootstrap.DescribeReturns(&model.OpSpec{Inputs: map[string]*model.ParamSpec{"HOME": {Dir: &model.DirParamSpec{}}}}, nil)
		bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{}}, nil)
		caller := _opCaller{bootstrapper: bootstrap, caller: new(fakes.FakeCaller), dataDirPath: scratch}
		if _, err := caller.Call(ctx, call.Op, nil, "root", child.Op, scope); err != nil {
			t.Fatal(err)
		}
		if bootstrap.PlanCallCount() != 1 {
			t.Fatal("Go service inputs did not validate")
		}
		_, op, _, _ := bootstrap.PlanArgsForCall(0)
		if op.Inputs["HOME"] == nil || op.Inputs["HOME"].Dir == nil || *op.Inputs["HOME"].Dir != home {
			t.Fatal("caller HOME did not reach the Go service")
		}
		return
	}
	t.Fatal("missing Go service child op")
}
