package goconvert

import (
	"context"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"strings"
	"testing"
)

func TestComparePlans(t *testing.T) {
	source := []byte("name: compare\ninputs: {region: {string: {default: west}}}\noutputs: {result: {string: {}}, unused: {string: {default: none}}}\nrun:\n  serial:\n  - container: {image: {ref: alpine}, stdout: '$(result)'}\n  - container: {image: {ref: alpine}, envVars: {REGION: '$(region)', RESULT: '$(result)', COUNT: 2, FLAG: true, JSON: [a, b]}, cmd: [echo, 'prefix $(region)']}\n")
	a, err := Analyze(context.Background(), t.TempDir(), source, nil)
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := Render(a, RenderOptions{})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := goop.ParsePlan(runCode(t, rendered, "plan", map[string]any{"region": "west"}))
	if err != nil {
		t.Fatal(err)
	}
	desc := describeCode(t, rendered.Source)
	if err = Compare(a, desc, plan, map[string]any{"region": "west"}); err != nil {
		t.Fatal(err)
	}
	(*plan.Run.Serial)[1].Container.EnvVars.(map[string]any)["REGION"] = "east"
	err = Compare(a, desc, plan, map[string]any{"region": "west"})
	if err == nil || !strings.Contains(err.Error(), "envVars.REGION") {
		t.Fatalf("wrong diff: %v", err)
	}
	desc.Description = "changed"
	if err = Compare(a, desc, plan, nil); err == nil || !strings.Contains(err.Error(), "description") {
		t.Fatalf("wrong interface diff: %v", err)
	}
}

func TestVerificationInputs(t *testing.T) {
	spec, err := opfile.Unmarshal([]byte("name: inputs\ninputs:\n  deploy: {boolean: {default: false}}\n  region: {string: {constraints: {enum: [west, east]}}}\n  secret: {string: {isSecret: true, default: marker}}\nrun: {container: {image: {ref: alpine}, cmd: ['$(region)']}}\n"))
	if err != nil {
		t.Fatal(err)
	}
	sets, err := InputSets(spec, t.TempDir(), map[string]any{"secret": "marker"})
	if err != nil {
		t.Fatal(err)
	}
	if len(sets) < 2 || len(sets) > 16 || sets[0]["region"] != "west" || sets[0]["deploy"] != false {
		t.Fatal(sets)
	}
	for _, set := range sets {
		if _, ok := set["secret"]; ok {
			t.Fatal("secret leaked")
		}
	}
	spec.Inputs["region"] = &model.ParamSpec{String: &model.StringParamSpec{Constraints: model.Constraints{"pattern": "^unusual-[0-9]{8}$"}}}
	if _, err = InputSets(spec, t.TempDir(), nil); err == nil || !strings.Contains(err.Error(), "pass -a region=") {
		t.Fatal(err)
	}
	sets, err = InputSets(spec, t.TempDir(), map[string]any{"region": "unusual-12345678"})
	if err != nil || sets[0]["region"] != "unusual-12345678" {
		t.Fatalf("%v %v", sets, err)
	}
}

func TestCompareUnusedStubOutputs(t *testing.T) {
	a, err := Analyze(context.Background(), t.TempDir(), []byte("name: unused-stub-output\nrun: {op: {ref: ../child}}\n"), func(context.Context, string) (*model.OpSpec, error) {
		return &model.OpSpec{Outputs: map[string]*model.ParamSpec{"message": {String: &model.StringParamSpec{}}}}, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := Render(a, RenderOptions{ModuleRoot: a.Dir, ModulePath: "ops.local/test"})
	if err != nil {
		t.Fatal(err)
	}
	desc, err := goop.ParseDescribe(runCode(t, rendered, "describe", nil))
	if err != nil {
		t.Fatal(err)
	}
	plan, err := goop.ParsePlan(runCode(t, rendered, "plan", nil))
	if err != nil {
		t.Fatal(err)
	}
	if err = Compare(a, desc, plan, nil); err != nil {
		t.Fatal(err)
	}
}

func TestPredicateBoundaryCandidatesComeFirst(t *testing.T) {
	spec, err := opfile.Unmarshal([]byte("name: boundary\ninputs: {num: {number: {default: 0}}}\nrun: {serial: [{container: {image: {ref: alpine}}}, {if: [{gt: ['$(num)', 5]}], container: {image: {ref: alpine}}}]}\n"))
	if err != nil {
		t.Fatal(err)
	}
	sets, err := InputSets(spec, t.TempDir(), nil)
	if err != nil {
		t.Fatal(err)
	}
	reached := map[bool]bool{}
	for _, set := range sets {
		reached[set["num"].(float64) > 5] = true
	}
	if !reached[true] || !reached[false] {
		t.Fatalf("missed predicate side: %v", sets)
	}
	if sets[0]["num"] != float64(0) {
		t.Fatalf("lost default set: %v", sets[0])
	}
}

func TestControlInputsPrecedeOtherVariations(t *testing.T) {
	source := "name: late-input\ninputs:\n"
	for _, name := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		source += fmt.Sprintf("  %s: {string: {default: one, constraints: {enum: [one, two, three, four]}}}\n", name)
	}
	source += "  zone: {string: {default: dev}}\n  zones: {array: {default: []}}\nrun:\n  serial:\n  - container: {image: {ref: alpine}, cmd: ['$(a) $(b) $(c) $(d) $(e) $(f) $(g) $(h)']}\n  - if: [{eq: ['$(zone)', prod]}]\n    container: {image: {ref: alpine}}\n  - serialLoop: {range: '$(zones)', run: {container: {image: {ref: alpine}}}}\n"
	spec, err := opfile.Unmarshal([]byte(source))
	if err != nil {
		t.Fatal(err)
	}
	sets, err := InputSets(spec, t.TempDir(), nil)
	if err != nil {
		t.Fatal(err)
	}
	prod, loop, flipped := false, false, false
	for _, set := range sets {
		prod = prod || set["zone"] == "prod"
		loop = loop || len(set["zones"].([]any)) == 2
		all := set["zone"] != "dev" && len(set["zones"].([]any)) > 0
		for _, name := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
			all = all && set[name] != "one"
		}
		flipped = flipped || all
	}
	if len(sets) > 16 || !prod || !loop || !flipped {
		t.Fatalf("starved control input or all-flipped set: %v", sets)
	}
}
