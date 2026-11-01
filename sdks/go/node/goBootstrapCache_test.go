package node

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	runtimeFakes "github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

func cachedBootstrapFixture(t *testing.T) (bootstrapper, *goBootstrapCache, *model.OpCall, *atomic.Int32) {
	t.Helper()
	root := t.TempDir()
	op := filepath.Join(root, "op")
	data := filepath.Join(root, "data")
	if err := os.Mkdir(op, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(op, "go.mod"), []byte("module hello\ngo 1.21\n"), 0600); err != nil {
		t.Fatal(err)
	}
	runtime := new(runtimeFakes.FakeContainerRuntime)
	runtime.ResolveImageReturns(&containerruntime.ImageIdentity{Ref: "sha256:compiler", ID: "sha256:compiler", OS: "linux", Arch: "amd64"}, nil)
	cache := &goBootstrapCache{runtime: runtime}
	builds := new(atomic.Int32)
	caller := new(fakes.FakeCaller)
	caller.CallStub = func(ctx context.Context, _ string, _ map[string]*model.Value, spec *model.CallSpec, _ string, _ *string, _ string) (map[string]*model.Value, error) {
		switch *spec.Name {
		case "bootstrap build":
			builds.Add(1)
			path := *spec.Container.Dirs["/opctl/build"].(model.Value).Dir
			if err := os.WriteFile(filepath.Join(path, "program"), []byte("synthetic executable"), 0700); err != nil {
				return nil, err
			}
			return nil, os.WriteFile(filepath.Join(path, "go-version"), []byte("go version go1.21 linux/amd64"), 0600)
		case "bootstrap describe":
			_, err := stdoutCaptureFrom(ctx).Write([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello"}}`))
			return nil, err
		case "bootstrap plan":
			_, err := stdoutCaptureFrom(ctx).Write([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello","run":{"container":{"image":{"ref":"alpine:3"},"cmd":["echo","new plan"]}}}}`))
			return nil, err
		}
		return nil, errors.New("unexpected bootstrap call")
	}
	return newBootstrapper(caller, data, cache), cache, &model.OpCall{BaseCall: model.BaseCall{OpPath: op}}, builds
}
func TestCompiledBootstrapCache(t *testing.T) {
	bootstrap, cache, op, builds := cachedBootstrapFixture(t)
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
				t.Error(err)
			}
		}()
	}
	group.Wait()
	if builds.Load() != 1 {
		t.Fatalf("built %d times", builds.Load())
	}
	for i := 0; i < 2; i++ {
		if _, err := bootstrap.Plan(context.Background(), op, "", map[string]interface{}{"input": i}); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := os.ReadDir(cache.directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("entries: %d", len(entries))
	}
	path := filepath.Join(cache.directory, entries[0].Name())
	files, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 {
		t.Fatal("program entry retained extra files")
	}
	if err := os.WriteFile(filepath.Join(path, "program"), []byte("corrupt executable"), 0700); err != nil {
		t.Fatal(err)
	}
	if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 2 {
		t.Fatal("corrupt binary was reused")
	}
	// A restart discards memory and reuses the executable.
	restarted := &goBootstrapCache{runtime: cache.runtime}
	clone := bootstrap.(_bootstrapper)
	clone.cache = restarted
	if _, err := clone.Describe(context.Background(), op, ""); err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 2 {
		t.Fatal("restart rebuilt")
	}
	if err := os.WriteFile(filepath.Join(op.OpPath, "asset"), []byte("new source"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
		t.Fatal(err)
	}
	if builds.Load() != 3 {
		t.Fatal("source change reused program")
	}
}

func TestGitFetchedBootstrapCache(t *testing.T) {
	bootstrap, _, op, builds := cachedBootstrapFixture(t)
	data := bootstrap.(_bootstrapper).dataDirPath
	repository := filepath.Join(data, "ops", "example.com", "repo", "v1")
	fetchedOp := filepath.Join(repository, "ops", "hello")
	if err := os.MkdirAll(fetchedOp, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Rename(filepath.Join(op.OpPath, "go.mod"), filepath.Join(repository, "go.mod")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "go.mod"), []byte("module hello\ngo 1.21\nreplace example.com/helper => ./helper\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(repository, "helper"), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(repository, "helper", "go.mod"), []byte("module example.com/helper\ngo 1.21\n"), 0600); err != nil {
		t.Fatal(err)
	}
	op.OpPath, op.RepositoryRoot = fetchedOp, repository
	for i := 0; i < 2; i++ {
		if _, err := bootstrap.Plan(context.Background(), op, "", nil); err != nil {
			t.Fatal(err)
		}
	}
	if builds.Load() != 1 {
		t.Fatalf("fetched op built %d times", builds.Load())
	}
}
func TestBootstrapCacheWaiterCancellation(t *testing.T) {
	_, cache, _, _ := cachedBootstrapFixture(t)
	owner, _, err := cache.acquire(context.Background(), "key")
	if err != nil {
		t.Fatal(err)
	}
	defer owner()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := cache.acquire(ctx, "key"); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled waiter: %v", err)
	}
	other, _, err := cache.acquire(context.Background(), "different")
	if err != nil {
		t.Fatal(err)
	}
	other()
}

func TestBootstrapCacheReportsWaiting(t *testing.T) {
	cache := &goBootstrapCache{}
	release, waited, err := cache.acquire(context.Background(), "key")
	if err != nil || waited {
		t.Fatalf("uncontended acquire: waited=%v, err=%v", waited, err)
	}
	defer release()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, waited, err = cache.acquire(ctx, "key")
	if !waited || !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("contended acquire: waited=%v, err=%v", waited, err)
	}
}
func TestBootstrapProgramLeasesAndOversize(t *testing.T) {
	bootstrap, cache, op, builds := cachedBootstrapFixture(t)
	cache.budget = 1
	if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
		t.Fatal(err)
	}
	entries, err := os.ReadDir(cache.directory)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatal("oversized executable persisted after release")
	}
	cache.budget = 1 << 20
	program, release, err := cache.prepare(context.Background(), bootstrap.(_bootstrapper), op, "")
	if err != nil {
		t.Fatal(err)
	}
	cache.budget = 1
	if err := cache.prune(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(program.path, "program")); err != nil {
		t.Fatal("active executable pruned")
	}
	release()
	if err := cache.prune(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(program.path); !os.IsNotExist(err) {
		t.Fatal("unused executable retained beyond budget")
	}
	if builds.Load() != 2 {
		t.Fatal("unexpected build count")
	}
	// Abandoned staging is removed lazily at node startup.
	stage := filepath.Join(cache.directory, ".tmp-abandoned")
	if err := os.Mkdir(stage, 0700); err != nil {
		t.Fatal(err)
	}
	fresh := &goBootstrapCache{}
	if err := fresh.initialize(bootstrap.(_bootstrapper).dataDirPath); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(stage); !os.IsNotExist(err) {
		t.Fatal("abandoned staging retained")
	}
}

func TestBootstrapBuildFailuresAndSourceChanges(t *testing.T) {
	for _, changes := range []int{0, 1, 2} {
		bootstrap, cache, op, builds := cachedBootstrapFixture(t)
		caller := bootstrap.(_bootstrapper).caller.(*fakes.FakeCaller)
		original := caller.CallStub
		attempts := 0
		caller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, root string) (map[string]*model.Value, error) {
			if *spec.Name == "bootstrap build" {
				attempts++
				if changes == 0 && attempts == 1 {
					return nil, errors.New("build failed")
				}
			}
			out, err := original(ctx, id, scope, spec, path, parent, root)
			if *spec.Name == "bootstrap build" && attempts <= changes {
				err = os.WriteFile(filepath.Join(op.OpPath, "asset"), []byte{byte(attempts)}, 0600)
			}
			return out, err
		}
		_, err := bootstrap.Describe(context.Background(), op, "")
		if changes == 1 && err != nil {
			t.Fatal(err)
		}
		if changes != 1 && err == nil {
			t.Fatal("expected preparation failure")
		}
		entries, readErr := os.ReadDir(cache.directory)
		if readErr != nil {
			t.Fatal(readErr)
		}
		if changes != 1 && len(entries) != 0 {
			t.Fatal("failure published a cache entry")
		}
		if changes == 0 {
			if _, err := bootstrap.Describe(context.Background(), op, ""); err != nil {
				t.Fatal(err)
			}
		}
		if changes == 1 && builds.Load() != 2 {
			t.Fatal("source edit did not rebuild")
		}
	}
}

func TestCancelledBuildOwnerDoesNotCancelWaiter(t *testing.T) {
	bootstrap, _, op, builds := cachedBootstrapFixture(t)
	caller := bootstrap.(_bootstrapper).caller.(*fakes.FakeCaller)
	original := caller.CallStub
	var first atomic.Bool
	entered := make(chan struct{})
	caller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, root string) (map[string]*model.Value, error) {
		if *spec.Name == "bootstrap build" && first.CompareAndSwap(false, true) {
			close(entered)
			<-ctx.Done()
			return nil, ctx.Err()
		}
		return original(ctx, id, scope, spec, path, parent, root)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	owner := make(chan error, 1)
	go func() { _, err := bootstrap.Describe(ctx, op, ""); owner <- err }()
	<-entered
	waiter := make(chan error, 1)
	go func() { _, err := bootstrap.Describe(context.Background(), op, ""); waiter <- err }()
	cancel()
	if err := <-owner; !errors.Is(err, context.Canceled) {
		t.Fatalf("owner: %v", err)
	}
	if err := <-waiter; err != nil {
		t.Fatalf("waiter: %v", err)
	}
	if builds.Load() != 1 {
		t.Fatal("waiter did not build its executable")
	}
}
