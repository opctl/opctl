package node

import (
	"context"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanOp(t *testing.T) {
	root := t.TempDir()
	data := []byte("name: yaml\ninputs:\n  env:\n    string: {default: dev}\n  token:\n    string: {isSecret: true, default: hidden}\nrun:\n  serial:\n    - if: [{eq: ['$(env)', prod]}]\n      container:\n        image: {ref: alpine}\n    - container:\n        image: {ref: alpine}\n        envVars: {TOKEN: '$(token)'}\n")
	if err := os.WriteFile(filepath.Join(root, "op.yml"), data, 0600); err != nil {
		t.Fatal(err)
	}
	bootstrap := new(fakes.FakeBootstrapper)
	c := core{bootstrapper: bootstrap}
	req := model.PlanOpReq{Op: model.StartOpReqOp{Ref: root, PullCreds: &model.Creds{}}}
	plan, err := c.PlanOp(context.Background(), req)
	if err != nil || plan.Run.Container == nil || bootstrap.PlanCallCount() != 0 {
		t.Fatalf("%+v %v", plan, err)
	}
	secret := "SECRET-MARKER"
	req.Args = map[string]*model.Value{"token": {String: &secret}}
	_, err = c.PlanOp(context.Background(), req)
	if err == nil || !strings.Contains(err.Error(), "input token isn't passed to plan") || strings.Contains(err.Error(), secret) {
		t.Fatalf("%v", err)
	}
	if err = os.Remove(filepath.Join(root, "op.yml")); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(root, "op.go"), []byte("package main"), 0600); err != nil {
		t.Fatal(err)
	}
	spec := &model.OpSpec{Name: "go", Inputs: map[string]*model.ParamSpec{"env": {String: &model.StringParamSpec{Default: "dev"}}, "token": {String: &model.StringParamSpec{IsSecret: true}}}}
	bootstrap.DescribeReturns(spec, nil)
	bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}}}}}, nil)
	req.Args = nil
	plan, err = c.PlanOp(context.Background(), req)
	if err != nil || plan.Name != "go" || plan.Run.Container == nil {
		t.Fatalf("%+v %v", plan, err)
	}
	_, _, _, args := bootstrap.PlanArgsForCall(0)
	if len(args) != 1 || args["env"] != "dev" {
		t.Fatalf("sent private or missing inputs: %+v", args)
	}
}
