package node

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
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
