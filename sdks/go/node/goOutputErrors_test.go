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

func TestGoOutputBindingErrorKeepsPublicCause(t *testing.T) {
	for _, secret := range []bool{false, true} {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "op.go"), []byte("package main"), 0600); err != nil {
			t.Fatal(err)
		}
		bootstrap := newGoBootstrap()
		bootstrap.DescribeReturns(&model.OpSpec{Outputs: map[string]*model.ParamSpec{"out": {String: &model.StringParamSpec{IsSecret: secret}}}}, nil)
		bootstrap.PlanReturns(&goop.Plan{OpSpec: &model.OpSpec{Run: &model.CallSpec{Container: &model.ContainerCallSpec{}}}, OutputBindings: map[string]string{"out": "$(missing)"}}, nil)
		_, err := (_opCaller{caller: new(fakes.FakeCaller), bootstrapper: bootstrap, dataDirPath: dir}).Call(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: dir}, OpID: "root"}, nil, "root", &model.OpCallSpec{}, nil)
		if err == nil || !strings.Contains(err.Error(), "unable to bind Go output") {
			t.Fatal("missing Go binding context")
		}
		if strings.Contains(err.Error(), "missing") == secret {
			t.Fatal("public cause lost or sensitive cause exposed")
		}
	}
}
