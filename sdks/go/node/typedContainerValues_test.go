package node

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	nodeFakes "github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

func TestCollectedContainerDeclarations(t *testing.T) {
	for _, tc := range []struct {
		name, data, failure string
		declarations        map[string]*model.ParamSpec
	}{
		{"renamed", `{"count":3,"stdout":7}`, "", map[string]*model.ParamSpec{"count": {Number: &model.NumberParamSpec{}}, "stdout": {Number: &model.NumberParamSpec{}}}},
		{"wrong kind", `{"count":"3"}`, `expected number, received string`, map[string]*model.ParamSpec{"count": {Number: &model.NumberParamSpec{}}}},
		{"missing unbound", `{"count":3}`, `"unbound": required output is missing`, map[string]*model.ParamSpec{"unbound": {String: &model.StringParamSpec{}}}},
		{"wrong unbound", `{"count":3,"unbound":false}`, `"unbound": expected string, received boolean`, map[string]*model.ParamSpec{"unbound": {String: &model.StringParamSpec{}}}},
		{"constraint", `{"count":0}`, `count`, map[string]*model.ParamSpec{"count": {Number: &model.NumberParamSpec{Constraints: model.Constraints{"minimum": 1}}}}},
		{"secret constraint", `{"count":3,"unbound":{"key":"private-invalid-value"}}`, `sensitive`, map[string]*model.ParamSpec{"unbound": {Object: &model.ObjectParamSpec{IsSecret: true, Constraints: model.Constraints{"properties": map[string]interface{}{"key": map[string]interface{}{"enum": []interface{}{"allowed"}}}}}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			output := filepath.Join(dir, "outputs", "values.json")
			runtime := new(fakes.FakeContainerRuntime)
			// The runtime writes only to the interpreter-prepared scratch directory.
			runtime.RunContainerStub = func(_ context.Context, _ *model.ContainerCall, _ string, _ pubsub.EventPublisher, stdout, stderr io.WriteCloser) (*int64, error) {
				defer stdout.Close()
				defer stderr.Close()
				io.WriteString(stdout, "3\n")
				return nil, os.WriteFile(output, []byte(tc.data), 0600)
			}
			call := &model.ContainerCall{Image: &model.ContainerCallImage{Prepared: true}, OutputsFile: output}
			spec := &model.ContainerCallSpec{Stdout: "$(stream)", Outputs: map[string]string{"count": "$(renamed)"}}
			if tc.name == "renamed" {
				spec.Outputs["stdout"] = "$(member)"
			}
			caller := _containerCaller{containerRuntime: runtime, pubSub: new(nodeFakes.FakePubSub)}
			result, err := caller.callWithResults(context.Background(), call, nil, spec, "root", tc.declarations)
			if tc.failure != "" {
				if err == nil || !strings.Contains(err.Error(), tc.failure) {
					t.Fatalf("expected %q, received %v", tc.failure, err)
				}
				if len(result) != 0 {
					t.Fatal("failed declaration published partial outputs")
				}
				if strings.Contains(err.Error(), "private-invalid-value") {
					t.Fatal("sensitive constraint error contains data")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				if result["renamed"].Number == nil || *result["renamed"].Number != 3 || result["stream"].String == nil || *result["stream"].String != "3" || result["member"].Number == nil || *result["member"].Number != 7 {
					t.Fatal("boxed kinds or renamed values changed")
				}
			}
			if _, err := os.Stat(filepath.Dir(output)); !os.IsNotExist(err) {
				t.Fatal("outputs scratch was retained")
			}
		})
	}
}
