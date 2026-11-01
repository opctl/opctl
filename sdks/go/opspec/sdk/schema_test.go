package sdk

import (
	"strings"
	"testing"
)

func TestSchemaSelection(t *testing.T) {
	graphs := []*callSpec{
		{Op: &opCallSpec{Ref: "opaque"}},
		{Container: &containerSpec{Stdout: "$(s)"}},
		{Container: &containerSpec{Outputs: map[string]string{"n": "$(n)"}}},
		{Container: &containerSpec{OutputParams: map[string]paramSpec{"n": {"number": map[string]any{}}}, Outputs: map[string]string{"n": "$(n)"}}},
	}
	for i, g := range graphs {
		for _, r := range []schemaRange{{1, 1}, {1, 2}, {1, 3}, {1, 4}, {2, 3}, {1, 99}, {5, 6}, {0, 1}, {3, 2}} {
			got, err := selectSchema(r, &callSpec{Serial: []*callSpec{{Parallel: []*callSpec{g}}}})
			high := r.Max
			if high > 4 {
				high = 4
			}
			valid := r.Min >= 1 && r.Max >= r.Min && r.Min <= 4 && i+1 <= high
			if (err == nil) != valid || valid && got != high {
				t.Fatalf("feature %d range %v: %d %v", i+1, r, got, err)
			}
			if !valid && r.Min == 1 && r.Max < i+1 && (!strings.Contains(err.Error(), "upgrade opctl") || !strings.Contains(err.Error(), "run.serial[0].parallel[0].container")) {
				t.Fatalf("unsafe or unhelpful diagnostic: %v", err)
			}
			if r.Min > 4 && !strings.Contains(err.Error(), "upgrade the Go op SDK") {
				t.Fatal(err)
			}
		}
	}
	// Names in ordinary op inputs/outputs do not consume schema versions.
	g := &callSpec{Op: &opCallSpec{Ref: "opaque", Inputs: map[string]any{"stdout": map[string]any{"outputParams": true}}, Outputs: map[string]string{"count": "$(n)"}}}
	if v, err := selectSchema(schemaRange{1, 1}, g); err != nil || v != 1 {
		t.Fatalf("literal misidentified: %v", err)
	}
}
