package goop

import (
	"os"
	"testing"
)

func TestTypedContainerOutputSharedFixture(t *testing.T) {
	data, err := os.ReadFile("../sdk/testdata/plan.typed-container-outputs.json")
	if err != nil {
		t.Fatal(err)
	}
	plan, err := ParsePlan(data)
	if err != nil {
		t.Fatal(err)
	}
	c := (*plan.Run.Serial)[0].Container
	if len(c.OutputParams) != 2 || c.OutputParams["count"].Number == nil || c.OutputParams["count"].Number.Constraints["minimum"] != float64(1) || c.Outputs["count"] != "$(opctl_number_1)" || c.Outputs["text"] != "$(opctl_string_2)" {
		t.Fatal("typed declarations or bindings changed at host boundary")
	}
}
