package node

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/internal/fakes"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

func TestBootstrapSchemaRangeAndMemo(t *testing.T) {
	b, initial, op, builds := cachedBootstrapFixture(t)
	original := b.(_bootstrapper)
	fake := original.caller.(*fakes.FakeCaller)
	build := fake.CallStub
	var describes, plans atomic.Int32
	fake.CallStub = func(ctx context.Context, id string, scope map[string]*model.Value, spec *model.CallSpec, path string, parent *string, root string) (map[string]*model.Value, error) {
		if *spec.Name == "bootstrap build" {
			return build(ctx, id, scope, spec, path, parent, root)
		}
		data, err := os.ReadFile(*scope["request"].File)
		if err != nil {
			return nil, err
		}
		var req goop.Request
		if err := json.Unmarshal(data, &req); err != nil {
			return nil, err
		}
		if req.OpspecSchema.Min != 1 {
			return nil, fmt.Errorf("unexpected minimum")
		}
		if *spec.Name == "bootstrap describe" {
			describes.Add(1)
			_, err = stdoutCaptureFrom(ctx).Write([]byte(fmt.Sprintf(`{"contractVersion":1,"opspecSchema":%d,"opspec":{"name":"hello","inputs":{"text":{"string":{"default":"original"}}}}}`, req.OpspecSchema.Max)))
			return nil, err
		}
		plans.Add(1)
		if req.OpspecSchema.Max < 4 {
			return nil, fmt.Errorf("this op uses typed container outputs (opspec >= 4); upgrade opctl to a version supporting opspec 4")
		}
		_, err = stdoutCaptureFrom(ctx).Write([]byte(`{"contractVersion":1,"opspecSchema":4,"opspec":{"name":"hello","run":{"container":{"image":{"ref":"alpine"}}}}}`))
		return nil, err
	}
	for _, max := range []int{1, 3, 4} {
		cache := &goBootstrapCache{runtime: initial.runtime, schema: goop.SchemaRange{Min: 1, Max: max}}
		bootstrap := original
		bootstrap.cache = cache
		for i := 0; i < 2; i++ {
			spec, err := bootstrap.Describe(context.Background(), op, "")
			if err != nil {
				t.Fatal(err)
			}
			if spec.Inputs["text"].String.Default != "original" {
				t.Fatal("memo shares parsed maps")
			}
			spec.Inputs["text"].String.Default = "mutated"
			_, err = bootstrap.Plan(context.Background(), op, "", map[string]interface{}{"input": i})
			if (err == nil) != (max == 4) {
				t.Fatalf("max %d: %v", max, err)
			}
			if max < 4 && !strings.Contains(err.Error(), "upgrade opctl") {
				t.Fatal(err)
			}
		}
		for key := range cache.describes {
			if !strings.Contains(key, fmt.Sprintf(":1:1:%d", max)) {
				t.Fatalf("describe memo key omits contract/range: %s", key)
			}
		}
	}
	if builds.Load() != 1 || describes.Load() != 3 || plans.Load() != 6 {
		t.Fatalf("builds=%d describes=%d plans=%d", builds.Load(), describes.Load(), plans.Load())
	}
	entries, err := os.ReadDir(filepath.Join(original.dataDirPath, "dcg"))
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatal("bootstrap request scratch retained")
	}
}

func TestNodeRejectsInvalidSchemaRangeBeforeCreatingState(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "unused")
	if _, err := newWithSchemaRange(context.Background(), nil, dir, goop.SchemaRange{Min: 1, Max: 5}); err == nil {
		t.Fatal("nonexistent reader advertised")
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("invalid range created node state")
	}
}
