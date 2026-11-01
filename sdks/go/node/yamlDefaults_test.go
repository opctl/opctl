package node

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

func TestYAMLSecretDefaultsKeepLegacyEvents(t *testing.T) {
	for _, kind := range []string{"inputs", "outputs"} {
		t.Run(kind, func(t *testing.T) {
			dir := t.TempDir()
			declaration := "name: defaults\n" + kind + ":\n  token:\n    string: {isSecret: true, default: legacy-token}\nrun:\n  container:\n    image: {ref: alpine}\n"
			if err := os.WriteFile(filepath.Join(dir, "op.yml"), []byte(declaration), 0600); err != nil {
				t.Fatal(err)
			}
			child := new(fakes.FakeCaller)
			child.CallStub = func(_ context.Context, _ string, scope map[string]*model.Value, _ *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
				if kind == "inputs" {
					value := scope["token"]
					if value == nil || value.Sensitive || value.String == nil || *value.String != "legacy-token" {
						t.Errorf("YAML defaulted input changed: %+v", value)
					}
				}
				return nil, nil
			}
			bus := new(fakes.FakePubSub)
			events := make(chan model.Event)
			close(events)
			bus.SubscribeReturns(events, nil)
			caller := _caller{pubSub: bus, opCaller: _opCaller{caller: child, dataDirPath: dir}, dataDirPath: dir}
			bindings := map[string]string{}
			if kind == "outputs" {
				bindings["token"] = "$(renamed)"
			}
			_, err := caller.Call(context.Background(), "root", nil, &model.CallSpec{Op: &model.OpCallSpec{Ref: dir, PullCreds: &model.CredsSpec{}, Outputs: bindings}}, dir, nil, "root")
			if err != nil {
				t.Fatal(err)
			}
			if kind == "outputs" {
				end := bus.PublishArgsForCall(bus.PublishCallCount() - 1).CallEnded
				if end == nil || end.Outputs["renamed"] == nil || end.Outputs["renamed"].String == nil || *end.Outputs["renamed"].String != "legacy-token" {
					t.Fatalf("CallEnded dropped the YAML defaulted secret output: %+v", end)
				}
			} else {
				start := bus.PublishArgsForCall(0).CallStarted
				if start == nil || start.Call.Op.Inputs["token"] == nil || *start.Call.Op.Inputs["token"].String != "legacy-token" {
					t.Fatalf("CallStarted dropped the YAML defaulted secret input: %+v", start)
				}
			}
		})
	}
}

func TestYAMLSecretDefaultErrorsKeepLegacyMessages(t *testing.T) {
	for _, declaration := range []string{
		"number: {default: legacy-invalid-number}",
		"string: {default: '$(missing)'}",
		"string: {default: legacy-invalid-enum, constraints: {enum: [allowed]}}",
	} {
		t.Run(declaration, func(t *testing.T) {
			var publicError string
			for _, secret := range []bool{false, true} {
				dir := t.TempDir()
				param := declaration
				if secret {
					param = strings.Replace(param, "{", "{isSecret: true, ", 1)
				}
				data := "name: defaults\ninputs:\n  token:\n    " + param + "\nrun:\n  container:\n    image: {ref: alpine}\n"
				if err := os.WriteFile(filepath.Join(dir, "op.yml"), []byte(data), 0600); err != nil {
					t.Fatal(err)
				}
				bus := new(fakes.FakePubSub)
				caller := _caller{pubSub: bus, dataDirPath: dir}
				_, err := caller.Call(context.Background(), "root", nil, &model.CallSpec{Op: &model.OpCallSpec{Ref: dir, PullCreds: &model.CredsSpec{}}}, dir, nil, "root")
				if err == nil {
					t.Fatal("invalid default accepted")
				}
				message := strings.ReplaceAll(err.Error(), dir, "OP")
				if !secret {
					publicError = message
				} else if message != publicError {
					t.Errorf("secret default error changed:\n%s\nlegacy:\n%s", message, publicError)
				}
				end := bus.PublishArgsForCall(bus.PublishCallCount() - 1).CallEnded
				if end == nil || end.Error == nil || end.Error.Message != err.Error() {
					t.Fatal("CallEnded changed the input error message")
				}
			}
		})
	}
}
