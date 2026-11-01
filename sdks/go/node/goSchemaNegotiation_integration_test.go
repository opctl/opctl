package node

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

// Opt in on a Docker host. All ranges use the same real node and SDK implementation.
func TestGoSchemaNegotiationDocker(t *testing.T) {
	repo := os.Getenv("OPCTL_SCHEMA_NEGOTIATION_REPO")
	if repo == "" {
		t.Skip("set OPCTL_SCHEMA_NEGOTIATION_REPO on a Docker host")
	}
	for _, tc := range []struct {
		name string
		max  int
		ops  []string
	}{
		{"baseline-old", 1, []string{"baseline", "baseline", "yaml-baseline"}},
		{"typed-old", 3, []string{"typed", "typed", "yaml-typed"}},
		{"current", 4, []string{"baseline", "typed", "yaml-typed"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			data := t.TempDir()
			runtime, err := docker.New(ctx, "unix:///var/run/docker.sock")
			if err != nil {
				cancel()
				t.Fatal(err)
			}
			n, err := newWithSchemaRange(ctx, runtime, data, goop.SchemaRange{Min: 1, Max: tc.max})
			if err != nil {
				cancel()
				t.Fatal(err)
			}
			defer func() {
				cancel()
				if err := n.(core).stateStore.(*_stateStore).db.Close(); err != nil {
					t.Error(err)
				}
			}()
			for _, name := range tc.ops {
				path := filepath.Join(repo, "test-suite/run/go/schema-negotiation", name)
				// Describe both cold and warm through the normal Core API before StartOp.
				for i := 0; i < 2; i++ {
					if _, err := n.DescribeOp(ctx, model.DescribeOpReq{Op: model.StartOpReqOp{Ref: path}}); err != nil {
						t.Fatal(err)
					}
				}
				id, err := n.StartOp(ctx, model.StartOpReq{Op: model.StartOpReqOp{Ref: path}})
				if err != nil {
					t.Fatal(err)
				}
				events, err := n.GetEventStream(ctx, &model.GetEventStreamReq{Filter: model.EventFilter{Roots: []string{id}}})
				if err != nil {
					t.Fatal(err)
				}
				var containers, cachedDescribes int
				var root *model.CallEnded
			collect:
				for {
					select {
					case <-ctx.Done():
						t.Fatal("timed out waiting for root completion")
					case event, ok := <-events:
						if !ok {
							t.Fatal("event stream closed before root completion")
						}
						if start := event.CallStarted; start != nil && start.Call.Name != nil && *start.Call.Name == "bootstrap describe (cached)" {
							cachedDescribes++
						}
						if start := event.CallStarted; start != nil && start.Call.Container != nil && (start.Call.Name == nil || !strings.HasPrefix(*start.Call.Name, "bootstrap ")) {
							containers++
						}
						if end := event.CallEnded; end != nil && end.Call.ID == id {
							root = end
							break collect
						}
					}
				}
				if cachedDescribes == 0 {
					t.Fatal("warm describe did not use the node memo")
				}
				failed := tc.max < 4 && strings.Contains(name, "typed")
				if failed {
					want := "this op uses typed container outputs (opspec >= 4); upgrade opctl to a version supporting opspec 4"
					if root.Outcome != model.OpOutcomeFailed || root.Error == nil || !strings.Contains(root.Error.Message, want) {
						t.Fatalf("missing upgrade failure: outcome=%s error=%v", root.Outcome, root.Error)
					}
					if containers != 0 || len(root.Outputs) != 0 {
						t.Fatalf("incompatible plan started %d user containers or published results", containers)
					}
					t.Logf("%s [1,%d]: %s; user containers=0, outputs=0", name, tc.max, want)
				} else {
					if root.Outcome != model.OpOutcomeSucceeded || containers != 2 {
						t.Fatalf("%s: outcome=%s containers=%d error=%v", name, root.Outcome, containers, root.Error)
					}
					if strings.Contains(name, "typed") && (root.Outputs["count"] == nil || root.Outputs["count"].Number == nil || *root.Outputs["count"].Number != 3) {
						t.Fatal("typed Number result missing")
					}
					t.Logf("%s [1,%d]: succeeded; user containers=2", name, tc.max)
				}
				if err := filepath.WalkDir(filepath.Join(data, "dcg"), func(path string, entry os.DirEntry, err error) error {
					if err == nil && entry.Name() == "request.json" {
						t.Errorf("bootstrap request scratch retained")
					}
					return err
				}); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
