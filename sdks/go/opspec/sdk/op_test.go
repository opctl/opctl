package sdk

import (
	"encoding/json"
	"math"
	"strings"
	"testing"
)

func TestCallOpEncodesInputs(t *testing.T) {
	type Inputs struct {
		Label String `opctl:"label"`
	}
	values := map[string]any{"nested": map[string]any{"$key": []any{"$(literal)", 3, true}}, "empty": ""}
	base := CallOp("github.com/example/repo#v1/path").Inputs(values)
	values["empty"] = "changed"
	values["nested"].(map[string]any)["$key"] = "changed"
	data, err := Plan(func(in Inputs) Call {
		return Serial(base, CallOp("./yaml").Inputs(map[string]any{"label": in.Label}).PullCreds(" ", in.Label))
	}, map[string]any{"label": "public"})
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	for _, expected := range []string{`github.com/example/repo#v1/path`, `opctl_empty_string`, `\\$key`, `\\$(literal)`, `"label":"$(label)"`, `"password":"$(label)"`} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %s in %s", expected, data)
		}
	}
	cleared, err := Plan(func(struct{}) Call { return base.Inputs(nil) }, nil)
	if err != nil || strings.Contains(string(cleared), `"inputs"`) {
		t.Fatalf("clear inputs: %s %v", cleared, err)
	}
}
func TestCallOpRejectsInvalidInputs(t *testing.T) {
	cycle := map[string]any{}
	cycle["self"] = cycle
	for _, value := range []any{nil, math.NaN(), func() {}, make(chan int), cycle, map[int]string{1: "x"}, (*String)(nil)} {
		if _, err := Plan(func(struct{}) Call { return CallOp("./yaml").Inputs(map[string]any{"bad": value}).Inputs(nil) }, nil); err == nil {
			t.Fatalf("accepted %T", value)
		}
	}
	for _, ref := range []string{"", "$(ref)"} {
		if _, err := Plan(func(struct{}) Call { return CallOp(ref) }, nil); err == nil {
			t.Fatalf("accepted ref %q", ref)
		}
	}
	if _, err := Plan(func(struct{}) Call { return CallOp("./yaml").PullCreds(3, "password") }, nil); err == nil {
		t.Fatal("accepted non-string creds")
	}
}
func TestCallOpGolden(t *testing.T) {
	helloDir(t)
	data, err := Plan(func(struct{}) Call { return CallOp("./yaml").Inputs(map[string]any{"label": "from-go"}) }, nil)
	if err != nil {
		t.Fatal(err)
	}
	var actual any
	if err := json.Unmarshal(data, &actual); err != nil {
		t.Fatal(err)
	}
	if string(data) != "{\"contractVersion\":1,\"opspecSchema\":4,\"opspec\":{\"name\":\"hello\",\"run\":{\"op\":{\"ref\":\"./yaml\",\"inputs\":{\"label\":\"from-go\"}}}}}\n" {
		t.Fatalf("%s", data)
	}
}
