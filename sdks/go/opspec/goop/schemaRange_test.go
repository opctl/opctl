package goop

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/opspec"
	"os"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

func TestBootstrapSchemaRanges(t *testing.T) {
	for _, r := range []SchemaRange{{1, 4}, {1, 1}, {2, 3}, {4, 4}} {
		req, err := NewRequestForRange("plan", nil, r)
		if err != nil || req.OpspecSchema != r {
			t.Fatalf("range %v: %v", r, err)
		}
		for v := 1; v <= 5; v++ {
			data := []byte(fmt.Sprintf(`{"contractVersion":1,"opspecSchema":%d,"opspec":{"name":"baseline","run":{"op":{"ref":"./opaque","inputs":{"stdout":{"outputParams":true}},"outputs":{"count":"$(n)"}}}}}`, v))
			_, err := ParsePlanForRange(data, r)
			if (err == nil) != (v >= r.Min && v <= r.Max) {
				t.Fatalf("range %v schema %d: %v", r, v, err)
			}
			data = []byte(fmt.Sprintf(`{"contractVersion":1,"opspecSchema":%d,"opspec":{"name":"baseline"}}`, v))
			_, err = ParseDescribeForRange(data, r)
			if (err == nil) != (v >= r.Min && v <= r.Max) {
				t.Fatalf("describe range %v schema %d: %v", r, v, err)
			}
		}
	}
	for _, r := range []SchemaRange{{}, {0, 1}, {-1, 2}, {3, 2}, {1, 5}, {5, 6}} {
		if _, err := NewRequestForRange("plan", nil, r); err == nil {
			t.Fatalf("invalid range %v accepted", r)
		}
		if _, err := ParsePlanForRange(nil, r); err == nil {
			t.Fatalf("invalid parse range %v accepted", r)
		}
	}
	if opfile.MinSchemaVersion != 1 {
		t.Fatal("old readers must remain supported")
	}
}

func TestBootstrapFeatureLabels(t *testing.T) {
	for _, tc := range []struct {
		field   string
		version int
	}{
		{`"stdout":"$(s)"`, 2}, {`"stderr":"$(s)"`, 2}, {`"outputs":{"n":"$(n)"}`, 3},
		{`"outputs":{"n":"$(n)"},"outputParams":{"n":{"number":{}}}`, 4},
	} {
		leaf := `{"container":{"image":{"ref":"alpine"},` + tc.field + `}}`
		for _, graph := range []string{leaf, `{"serial":[{"parallel":[` + leaf + `]}]}`, `{"serialLoop":{"range":[1],"run":` + leaf + `}}`, `{"parallelLoop":{"range":[1],"run":` + leaf + `}}`} {
			for v := 1; v <= 4; v++ {
				data := []byte(fmt.Sprintf(`{"contractVersion":1,"opspecSchema":%d,"opspec":{"name":"features","run":%s}}`, v, graph))
				_, err := ParsePlanForRange(data, SchemaRange{1, 4})
				if (err == nil) != (v >= tc.version) {
					t.Fatalf("schema %d field %s: %v", v, tc.field, err)
				}
				if v < tc.version && (err == nil || !strings.Contains(err.Error(), "schema")) {
					t.Fatalf("missing feature diagnostic: %v", err)
				}
			}
		}
	}
	for _, field := range []string{`"stdout":null`, `"stdout":""`, `"outputs":{}`, `"outputs":null`, `"outputParams":{}`, `"outputParams":null`} {
		data := []byte(`{"contractVersion":1,"opspecSchema":1,"opspec":{"name":"invalid","run":{"container":{"image":{"ref":"alpine"},` + field + `}}}}`)
		if _, err := ParsePlan(data); err == nil {
			t.Fatalf("invalid field accepted: %s", field)
		}
	}
}

func TestSharedNegotiatedSchemaFixtures(t *testing.T) {
	for version := 1; version <= 4; version++ {
		root := fmt.Sprintf("../sdk/testdata/schema-versions/v%d/", version)
		data, err := os.ReadFile(root + "request.json")
		if err != nil {
			t.Fatal(err)
		}
		var req Request
		if err := json.Unmarshal(data, &req); err != nil {
			t.Fatal(err)
		}
		data, err = os.ReadFile(root + "describe.json")
		if err != nil {
			t.Fatal(err)
		}
		if _, err := ParseDescribeForRange(data, req.OpspecSchema); err != nil {
			t.Fatalf("v%d describe: %v", version, err)
		}
		data, err = os.ReadFile(root + "plan.json")
		if err != nil {
			t.Fatal(err)
		}
		plan, err := ParsePlanForRange(data, req.OpspecSchema)
		if err != nil {
			t.Fatalf("v%d plan: %v", version, err)
		}
		if needed := opspec.MinimumSchema(plan.Run); needed != version {
			t.Fatalf("v%d fixture only exercises schema %d", version, needed)
		}
		if version > 1 {
			if _, err := ParsePlanForRange(data, SchemaRange{1, version - 1}); err == nil {
				t.Fatal("response outside request accepted")
			}
		}
	}
	// Contract-1 op output bindings remain readable at schema 1.
	data, err := os.ReadFile("../sdk/testdata/plan.outputs.json")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := ParsePlanForRange(data, SchemaRange{1, 1})
	if err != nil || plan.OutputBindings["echoed"] != "$(text)" {
		t.Fatalf("schema-1 output bindings: %v", err)
	}
}

func TestNullEnvelopeVersionsDoNotReuseEarlierFields(t *testing.T) {
	for _, versions := range []string{
		`"contractVersion":1,"contractVersion":null,"opspecSchema":1`,
		`"contractVersion":1,"opspecSchema":1,"opspecSchema":null`,
	} {
		data := []byte(`{` + versions + `,"opspec":{"name":"hello","run":{"container":{"image":{"ref":"alpine"}}}}}`)
		if _, err := ParsePlan(data); err == nil {
			t.Fatal("null envelope version reused an earlier integer")
		}
	}
}
