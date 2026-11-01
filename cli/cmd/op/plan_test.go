package op

import (
	"bytes"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/fakes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPlanCommandKeepsSecretsOut(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "op.go"), []byte("package main"), 0600); err != nil {
		t.Fatal(err)
	}
	api := new(fakes.FakeNode)
	spec := &model.OpSpec{Name: "dev", Inputs: map[string]*model.ParamSpec{"env": {String: &model.StringParamSpec{Default: "dev"}}, "token": {String: &model.StringParamSpec{IsSecret: true, Default: "SECRET-DEFAULT"}}, "src": {Dir: &model.DirParamSpec{}}}, Outputs: map[string]*model.ParamSpec{"result": {String: &model.StringParamSpec{}}}, Run: &model.CallSpec{Container: &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, Stdout: "$(opctl_string_1)"}}}
	api.DescribeOpReturns(spec, nil)
	api.PlanOpReturns(&model.OpPlan{OpSpec: spec, OutputBindings: map[string]string{"result": "$(opctl_string_1)"}}, nil)
	saved := node
	node = api
	defer func() { node = saved }()
	cmd := newPlanCmd(cliparamsatisfier.New(nil))
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{dir, "-a", "env=prod", "-a", "token=SECRET-MARKER", "-a", "src=/SECRET-PATH"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	for _, marker := range []string{"SECRET-MARKER", "SECRET-DEFAULT", "SECRET-PATH", "opctl_string_1"} {
		if strings.Contains(stdout.String()+stderr.String(), marker) {
			t.Fatalf("leaked %s: %s\n%s", marker, &stdout, &stderr)
		}
	}
	for _, want := range []string{"# opctl op plan:", `#   env: "prod"`, "has a secret default", "stdout: $(result)"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("missing %q: %s", want, &stdout)
		}
	}
	_, req := api.PlanOpArgsForCall(0)
	if len(req.Args) != 1 || *req.Args["env"].String != "prod" || api.StartOpCallCount() != 0 {
		t.Fatalf("%+v", req)
	}
	if spec.Inputs["token"].String.Default != "SECRET-DEFAULT" {
		t.Fatal("mutated source")
	}
	cmd = newPlanCmd(cliparamsatisfier.New(nil))
	stdout.Reset()
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{dir})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `#   env: "dev"`) {
		t.Fatalf("default missing from header: %s", &stdout)
	}
}
