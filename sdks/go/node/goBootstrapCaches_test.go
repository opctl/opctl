package node

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	runtimeFakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

func TestSharedBootstrapCaches(t *testing.T) {
	root := t.TempDir()
	op := filepath.Join(root, "op")
	if err := os.Mkdir(op, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(op, "go.mod"), []byte("module hello\ngo 1.21\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runtime := new(runtimeFakes.FakeContainerRuntime)
	runtime.ResolveImageReturns(&containerruntime.ImageIdentity{Ref: "sha256:toolchain", ID: "sha256:toolchain", OS: "linux", Arch: "amd64"}, nil)
	cache := &goBootstrapCache{runtime: runtime}
	data := filepath.Join(root, "data")
	caller := new(fakes.FakeCaller)
	var lastMod, lastBuild string
	caller.CallStub = func(ctx context.Context, _ string, scope map[string]*model.Value, spec *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
		if *spec.Name != "bootstrap build" {
			if _, ok := spec.Container.Dirs["/opctl/mod"]; ok {
				t.Fatal("phase mounted module cache")
			}
			if _, ok := spec.Container.Dirs["/opctl/cache"]; ok {
				t.Fatal("phase mounted build cache")
			}
			if !spec.Container.ReadOnlyMounts["/opctl/program"] {
				t.Fatal("program was writable")
			}
			_, err := stdoutCaptureFrom(ctx).Write([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello"}}`))
			return nil, err
		}

		env := spec.Container.EnvVars.(map[string]interface{})
		for key, want := range map[string]string{"GOWORK": "off", "GOENV": "off", "GOTOOLCHAIN": "local", "CGO_ENABLED": "0", "GOOS": "linux", "GOARCH": "amd64", "GOMODCACHE": "/opctl/mod", "GOCACHE": "/opctl/cache"} {
			if env[key] != want {
				t.Errorf("%s = %v", key, env[key])
			}
		}
		mod := *spec.Container.Dirs["/opctl/mod"].(model.Value).Dir
		build := *spec.Container.Dirs["/opctl/cache"].(model.Value).Dir
		if lastMod != "" && (lastMod != mod || lastBuild != build) {
			t.Error("cache paths changed across phases")
		}
		lastMod, lastBuild = mod, build
		if _, err := os.Stat(mod); err != nil {
			t.Fatal(err)
		}
		if !spec.Container.ReadOnlyMounts["/opctl/launcher"] || spec.Container.ReadOnlyMounts["/opctl/mod"] {
			t.Error("incorrect cache/launcher permissions")
		}
		launcher := scope["launcher"]
		if launcher == nil || launcher.File == nil {
			t.Fatal("missing scope-file launcher")
		}
		output := *spec.Container.Dirs["/opctl/build"].(model.Value).Dir
		if err := os.WriteFile(filepath.Join(output, "program"), []byte("synthetic executable"), 0700); err != nil {
			return nil, err
		}
		if err := os.WriteFile(filepath.Join(output, "go-version"), []byte("go version go1.21 linux/amd64"), 0600); err != nil {
			return nil, err
		}

		return nil, nil
	}
	bootstrap := newBootstrapper(caller, data, cache)
	for i := 0; i < 2; i++ {
		if _, err := bootstrap.Describe(context.Background(), &model.OpCall{BaseCall: model.BaseCall{OpPath: op}}, ""); err != nil {
			t.Fatal(err)
		}
	}
	if lastMod != filepath.Join(data, "go", "mod") || lastBuild != filepath.Join(data, "go", "build") {
		t.Fatal("cache is outside node data")
	}
}
