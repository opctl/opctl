package goop

import (
	"os"
	"testing"
)

func TestLiteralHandleSharedFixture(t *testing.T) {
	data, err := os.ReadFile("../sdk/testdata/plan.literal-handles.json")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := ParsePlan(data)
	if err != nil {
		t.Fatal(err)
	}
	c := (*plan.Run.Serial)[0].Op
	if c.Inputs["empty"] != "$(opctl_empty_string)" || c.Inputs["number"] != float64(2) {
		t.Fatalf("literal inputs changed: %#v", c.Inputs)
	}
}
