package node

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gitclient "github.com/go-git/go-git/v5/plumbing/transport/client"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/opctl/opctl/sdks/go/data"
	datafs "github.com/opctl/opctl/sdks/go/data/fs"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

func TestBootstrapStubsAndCache(t *testing.T) {
	b, cache, op, builds := cachedBootstrapFixture(t)
	write := func(path, text string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(text), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write(filepath.Join(op.OpPath, "op.go"), "package main\nimport _ \"hello/opctlgen/greet\" //opctl:op ../yaml\n")
	callee := filepath.Join(filepath.Dir(op.OpPath), "yaml/op.yml")
	yaml := "name: greet\ninputs:\n  token:\n    string: {isSecret: true, default: do-not-render}\nrun:\n  serial: []\n"
	write(callee, yaml)
	bootstrap := b.(_bootstrapper)
	fakeCaller := bootstrap.caller.(*fakes.FakeCaller)
	originalCall := fakeCaller.CallStub
	var mounted string
	fakeCaller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, root string) (map[string]*model.Value, error) {
		if *spec.Name == "bootstrap build" {
			mounted = *spec.Container.Dirs["/opctl/stubs"].(model.Value).Dir
			if !spec.Container.ReadOnlyMounts["/opctl/stubs"] || !strings.Contains(spec.Container.EnvVars.(map[string]interface{})["GOFLAGS"].(string), "-overlay=/opctl/stubs/overlay.json") {
				t.Error("overlay mount or flag lost")
			}
		}
		return originalCall(ctx, id, scope, spec, path, parent, root)
	}
	source, err := goop.Discover(op.OpPath, "")
	if err != nil {
		t.Fatal(err)
	}
	stubs, err := bootstrap.stubs(context.Background(), op, "root", source)
	if err != nil {
		t.Fatal(err)
	}
	if len(stubs.files) != 1 || strings.Contains(string(stubs.files["greet"]), "do-not-render") {
		t.Fatal("missing stub or secret leaked")
	}
	scratch, err := stubs.materialize(bootstrap.dataDirPath, source.ModuleRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(scratch)
	data, err := os.ReadFile(filepath.Join(scratch, "overlay.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "/opctl/stubs/greet.go") {
		t.Fatalf("%s", data)
	}
	if _, err := b.Describe(context.Background(), op, "root"); err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 1 {
		t.Fatal("no cold build")
	}
	if _, err := os.Stat(mounted); !os.IsNotExist(err) {
		t.Fatal("stub scratch retained after build", err)
	}
	write(filepath.Join(op.OpPath, "opctlgen/greet/stub.go"), "stale")
	write(filepath.Join(op.OpPath, "opctlgen/greet/extra.go"), "malicious")
	if _, err := b.Describe(context.Background(), op, "root"); err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 1 {
		t.Fatal("host stubs invalidated compiled program")
	}
	scratch2, err := stubs.materialize(bootstrap.dataDirPath, source.ModuleRoot)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(scratch2)
	data, err = os.ReadFile(filepath.Join(scratch2, "overlay.json"))
	if err != nil {
		t.Fatal(err)
	}
	var overlay struct{ Replace map[string]string }
	if err := json.Unmarshal(data, &overlay); err != nil {
		t.Fatal(err)
	}
	if value, ok := overlay.Replace[filepath.Join(source.ModuleRoot, "opctlgen/greet/extra.go")]; !ok || value != "" {
		t.Fatal("stale source not hidden")
	}
	write(callee, strings.Replace(yaml, "run:", "  optional:\n    string: {default: value}\nrun:", 1))
	if _, err := b.Describe(context.Background(), op, "root"); err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 2 {
		t.Fatalf("interface change did not rebuild: %d", builds.Load())
	}
	if len(cache.describes) != 2 {
		t.Fatalf("describe memo not keyed by stubs: %d", len(cache.describes))
	}
}

func TestBootstrapStubCycleAndHint(t *testing.T) {
	b, _, op, _ := cachedBootstrapFixture(t)
	if err := os.WriteFile(filepath.Join(op.OpPath, "op.go"), []byte("package main\nimport _ \"hello/opctlgen/self\" //opctl:op ./\n"), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := b.Describe(context.Background(), op, "root")
	if err == nil || !strings.Contains(err.Error(), "op stubs form a cycle") || !strings.Contains(err.Error(), "CallOp") {
		t.Fatalf("%v", err)
	}
	set := goStubSet{uses: []goop.StubUse{{Name: "self", Ref: "./"}}, files: map[string][]byte{"self": []byte("new")}}
	if err := os.MkdirAll(filepath.Join(op.OpPath, "opctlgen/self"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(op.OpPath, "opctlgen/self/stub.go"), []byte("old"), 0644); err != nil {
		t.Fatal(err)
	}
	hint := set.hint(op.OpPath, "undefined: opctl.StubCall")
	for _, want := range []string{"self (./)", "out of date", "Go op SDK >=", "opctl op gen"} {
		if !strings.Contains(hint, want) {
			t.Fatalf("%s", hint)
		}
	}
}

func TestStubFreeRecipeUnchanged(t *testing.T) {
	legacy := struct {
		Launcher []byte
		Command  []interface{}
		Env      map[string]interface{}
	}{[]byte("launcher"), []interface{}{"go"}, map[string]interface{}{"GOFLAGS": "old"}}
	want, _ := json.Marshal(legacy)
	got, _ := json.Marshal(goop.BuildRecipe{Launcher: legacy.Launcher, Command: legacy.Command, Env: legacy.Env})
	if string(got) != string(want) {
		t.Fatalf("recipe changed: %s", got)
	}
}

func TestSharedDescribeParentsGoCallee(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "op.go"), []byte("package main"), 0644); err != nil {
		t.Fatal(err)
	}
	dir, err := data.Resolve(context.Background(), root, datafs.New())
	if err != nil {
		t.Fatal(err)
	}
	fake := new(fakes.FakeBootstrapper)
	fake.DescribeReturns(&model.OpSpec{Name: "callee"}, nil)
	spec, err := describeOpDir(context.Background(), dir, nil, "caller", "root", fake, "")
	if err != nil || spec.Name != "callee" {
		t.Fatalf("%v %v", spec, err)
	}
	_, call, rootID := fake.DescribeArgsForCall(0)
	if call.OpID != "caller" || rootID != "root" {
		t.Fatalf("%+v %s", call, rootID)
	}
}

func TestBootstrapStubBuildHintsMatchDiagnostics(t *testing.T) {
	for _, tc := range []struct {
		diagnostic string
		changed    bool
	}{
		{"./op.go:7: undefined: misspelledHelper", false},
		{"go: module download failed", false},
		{"./op.go:7: unknown field Gone in struct literal of type greet.Inputs", false},
		{"./opctlgen/greet/stub.go:8: undefined: opctl.StubCall", true},
		{`./op.go:7: cannot use value as "hello/opctlgen/greet".Inputs`, true},
	} {
		t.Run(tc.diagnostic, func(t *testing.T) {
			b, _, op, _ := cachedBootstrapFixture(t)
			if err := os.WriteFile(filepath.Join(op.OpPath, "op.go"), []byte("package main\nimport _ \"hello/opctlgen/greet\" //opctl:op ./yaml\n"), 0644); err != nil {
				t.Fatal(err)
			}
			if err := os.Mkdir(filepath.Join(op.OpPath, "yaml"), 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(filepath.Join(op.OpPath, "yaml/op.yml"), []byte("name: greet\nrun: {serial: []}\n"), 0644); err != nil {
				t.Fatal(err)
			}
			caller := b.(_bootstrapper).caller.(*fakes.FakeCaller)
			caller.CallStub = func(ctx context.Context, _ string, _ map[string]*model.Value, _ *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
				_, _ = stdoutCaptureFrom(ctx).diagnostics.Write([]byte(tc.diagnostic))
				return nil, errors.New("build failed")
			}
			_, err := b.Describe(context.Background(), op, "root")
			if err == nil || !strings.Contains(err.Error(), tc.diagnostic) {
				t.Fatalf("%v", err)
			}
			if strings.Contains(err.Error(), "these ops no longer declare") != tc.changed {
				t.Fatalf("misleading build hint: %v", err)
			}
		})
	}
}

func TestBootstrapStubFetchAuthHints(t *testing.T) {
	for _, tc := range []struct {
		status int
		auth   error
	}{
		{http.StatusUnauthorized, model.ErrDataProviderAuthentication{}},
		{http.StatusForbidden, model.ErrDataProviderAuthorization{}},
		{http.StatusNotFound, nil},
		{http.StatusServiceUnavailable, nil},
	} {
		t.Run(http.StatusText(tc.status), func(t *testing.T) {
			server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(tc.status) }))
			defer server.Close()
			original := gitclient.Protocols["https"]
			gitclient.InstallProtocol("https", githttp.NewClient(server.Client()))
			t.Cleanup(func() { gitclient.InstallProtocol("https", original) })
			b, _, op, _ := cachedBootstrapFixture(t)
			ref := server.URL + "/repo#v1"
			if err := os.WriteFile(filepath.Join(op.OpPath, "op.go"), []byte("package main\nimport _ \"hello/opctlgen/greet\" //opctl:op "+ref+"\n"), 0644); err != nil {
				t.Fatal(err)
			}
			_, err := b.Describe(context.Background(), op, "root")
			if err == nil {
				t.Fatal("expected a fetch error")
			}
			if strings.Contains(err.Error(), "add credentials with opctl auth add") != (tc.auth != nil) {
				t.Fatalf("misleading fetch hint: %v", err)
			}
			if tc.auth != nil && !errors.Is(err, tc.auth) {
				t.Fatalf("lost auth cause: %v", err)
			}
		})
	}
}
