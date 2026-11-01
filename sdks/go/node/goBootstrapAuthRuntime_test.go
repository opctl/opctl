package node

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	runtimeFakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	containerInterpreter "github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
)

func TestBootstrapAuthRuntimeBuildCount(t *testing.T) {
	bootstrap, cache, op, _ := cachedBootstrapFixture(t)
	marker := "synthetic-runtime-auth-marker"
	cache.stateStore = &bootstrapSnapshotStore{auths: []model.Auth{{Resources: "example.com/team", Creds: model.Creds{Username: "user", Password: marker}}}}
	runtime := cache.runtime.(*runtimeFakes.FakeContainerRuntime)
	bus := new(fakes.FakePubSub)
	containerCaller := newContainerCaller(runtime, bus, nil)
	caller := bootstrap.(_bootstrapper).caller.(*fakes.FakeCaller)
	caller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, _ *string, root string) (map[string]*model.Value, error) {
		call, err := containerInterpreter.Interpret(scope, spec.Container, id, path, bootstrap.(_bootstrapper).dataDirPath)
		if err != nil {
			return nil, err
		}
		return containerCaller.Call(ctx, call, scope, spec.Container, root)
	}
	builds, describes, plans := 0, 0, 0
	runtime.RunContainerStub = func(ctx context.Context, call *model.ContainerCall, _ string, _ pubsub.EventPublisher, stdout, stderr io.WriteCloser) (*int64, error) {
		defer stdout.Close()
		defer stderr.Close()
		serialized, err := json.Marshal(call)
		if err != nil {
			return nil, err
		}
		if bytes.Contains(serialized, []byte(marker)) {
			t.Fatal("credentials entered runtime call config")
		}
		input := containerruntime.BootstrapInputFrom(ctx)
		if output := call.Dirs["/opctl/build"]; output != "" {
			builds++
			want := []string{"sh", "/opctl/launcher", "sh", "-ce", "go version > /opctl/build/go-version; go build -mod=readonly -buildvcs=false -trimpath -o /opctl/build/program ."}
			if !reflect.DeepEqual(call.Cmd, want) || input == nil {
				t.Fatal("build did not use the fixed compiler command with private input")
			}
			if err := os.WriteFile(filepath.Join(output, "program"), []byte("runtime-built executable"), 0700); err != nil {
				return nil, err
			}
			if err := os.WriteFile(filepath.Join(output, "go-version"), []byte("go version go1.21 linux/amd64"), 0600); err != nil {
				return nil, err
			}
		} else {
			if input != nil || len(call.Cmd) != 4 || call.Cmd[2] != "/opctl/program" || call.Dirs["/opctl/mod"] != "" || call.Dirs["/opctl/cache"] != "" {
				t.Fatal("cached phase received credentials/caches or invoked a compiler")
			}
			if call.Cmd[3] == "describe" {
				describes++
				_, err = io.WriteString(stdout, `{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello"}}`)
			} else {
				plans++
				_, err = io.WriteString(stdout, `{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello","run":{"container":{"image":{"ref":"alpine:3"}}}}}`)
			}
			if err != nil {
				return nil, err
			}
		}
		zero := int64(0)
		return &zero, nil
	}
	for i := 0; i < 2; i++ {
		if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
			t.Fatal(err)
		}
		if _, err := bootstrap.Plan(context.Background(), op, "", map[string]interface{}{"input": i}); err != nil {
			t.Fatal(err)
		}
	}
	if builds != 1 || describes != 1 || plans != 2 {
		t.Fatalf("runtime command counts: build=%d describe=%d plan=%d", builds, describes, plans)
	}
}
