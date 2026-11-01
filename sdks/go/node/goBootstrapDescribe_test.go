package node

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
)

func TestBootstrapDescribeMemo(t *testing.T) {
	bootstrap, cache, op, _ := cachedBootstrapFixture(t)
	bus := new(fakes.FakePubSub)
	cache.pubSub = bus
	caller := bootstrap.(_bootstrapper).caller.(*fakes.FakeCaller)
	original := caller.CallStub
	var describes atomic.Int32
	caller.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, root string) (map[string]*model.Value, error) {
		if *spec.Name == "bootstrap describe" {
			describes.Add(1)
			_, err := stdoutCaptureFrom(ctx).Write([]byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"hello","inputs":{"text":{"string":{"default":"original"}}},"outputs":{"echo":{"string":{}}}}}`))
			return nil, err
		}
		return original(ctx, id, scope, spec, path, parent, root)
	}
	var group sync.WaitGroup
	for i := 0; i < 8; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			spec, err := bootstrap.Describe(context.Background(), op, "")
			if err != nil {
				t.Error(err)
				return
			}
			spec.Inputs["text"].String.Default = "mutated"
			delete(spec.Outputs, "echo")
		}()
	}
	group.Wait()
	spec, err := bootstrap.Describe(context.Background(), op, "")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Inputs["text"].String.Default != "original" || spec.Outputs["echo"] == nil {
		t.Fatal("memo shared mutable declaration maps")
	}
	if describes.Load() != 1 {
		t.Fatalf("executed describe %d times", describes.Load())
	}
	if bus.PublishCallCount() != 16 {
		t.Fatalf("cached lifecycle events: %d", bus.PublishCallCount())
	}
	for i := 0; i < bus.PublishCallCount(); i++ {
		event := bus.PublishArgsForCall(i)
		if event.CallStarted != nil && event.CallStarted.Call.Container != nil {
			t.Fatal("memo hit invented a container")
		}
	}
}
func TestDescribeMemoBounds(t *testing.T) {
	cache := &goBootstrapCache{}
	for i := 0; i < 130; i++ {
		cache.rememberDescribe(string(rune(i)), []byte("valid response"))
	}
	if len(cache.describes) != 128 {
		t.Fatal("memo entry bound")
	}
	if _, ok := cache.cachedDescribe(string(rune(0))); ok {
		t.Fatal("LRU entry retained")
	}
	cache.rememberDescribe("too large", make([]byte, 64*1024*1024+1))
	if len(cache.describes) != 128 {
		t.Fatal("oversized response cached")
	}
}
