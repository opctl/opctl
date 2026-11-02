package snapshot

import (
	"encoding/json"
	"github.com/opctl/opctl/sdks/go/model"
	"strings"
	"testing"
)

func TestExpand(t *testing.T) {
	var spec model.OpSpec
	if err := json.Unmarshal([]byte(`{"name":"example","run":{"serial":[{"if":[{"eq":["$(env)","prod"]}],"container":{"image":{"ref":"alpine"},"cmd":["deploy"]}},{"serialLoop":{"range":"$(items)","vars":{"index":"$(i)","value":"$(v)"},"run":{"op":{"ref":"../child","inputs":{"item":"$(v)","message":"$(i):$(v)"}}}}}]}}`), &spec); err != nil {
		t.Fatal(err)
	}
	out, err := Expand(&spec, map[string]any{"env": "dev", "items": []any{"a", "$(literal)"}})
	if err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(out.Run)
	want := `{"serial":[{"op":{"ref":"../child","inputs":{"item":"a","message":"0:a"}}},{"op":{"ref":"../child","inputs":{"item":"\\$(literal)","message":"1:\\$(literal)"}}}]}`
	if string(data) != want {
		t.Fatalf("got %s; want %s", data, want)
	}
	original, _ := json.Marshal(spec.Run)
	if !strings.Contains(string(original), "serialLoop") {
		t.Fatal("mutated source")
	}
	out, err = Expand(&spec, map[string]any{"env": "prod", "items": []any{}})
	if err != nil || out.Run.Container == nil {
		t.Fatalf("%+v %v", out, err)
	}
}

func TestUnexpandableNodes(t *testing.T) {
	var spec model.OpSpec
	_ = json.Unmarshal([]byte(`{"name":"example","run":{"if":[{"exists":"$(file)"}],"serialLoop":{"range":"$(produced)","until":[{"eq":[1,1]}],"run":{"container":{"image":{"ref":"alpine"}}}}}}`), &spec)
	out, err := Expand(&spec, nil)
	if err != nil || out.Run.If == nil || out.Run.SerialLoop == nil {
		t.Fatalf("%+v %v", out, err)
	}
}

func TestShorthandOutputsShadowInputValues(t *testing.T) {
	for _, producer := range []string{`{"op":{"ref":"../child","outputs":{"version":""}}}`, `{"container":{"image":{"ref":"alpine"},"outputs":{"version":""}}}`} {
		var spec model.OpSpec
		if err := json.Unmarshal([]byte(`{"run":{"serial":[`+producer+`,{"if":[{"eq":["$(version)","old"]}],"container":{"image":{"ref":"alpine"}}}]}}`), &spec); err != nil {
			t.Fatal(err)
		}
		out, err := Expand(&spec, map[string]any{"version": "old"})
		if err != nil {
			t.Fatal(err)
		}
		if out.Run.Serial == nil || len(*out.Run.Serial) != 2 || (*out.Run.Serial)[1].If == nil {
			t.Fatalf("folded using shadowed input: %+v", out.Run)
		}
	}
}

func TestParallelLoopInvalidatesOuterValuesAfterAllIterations(t *testing.T) {
	var spec model.OpSpec
	source := `{"run":{"serial":[{"parallelLoop":{"range":[1,2],"run":{"if":[{"eq":["$(version)","old"]}],"op":{"ref":"../child","outputs":{"version":"$(version)"}}}}},{"if":[{"eq":["$(version)","old"]}],"container":{"image":{"ref":"alpine"}}}]}}`
	if err := json.Unmarshal([]byte(source), &spec); err != nil {
		t.Fatal(err)
	}
	out, err := Expand(&spec, map[string]any{"version": "old"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Run.Serial == nil || len(*out.Run.Serial) != 2 {
		t.Fatalf("lost calls: %+v", out.Run)
	}
	parallel := (*out.Run.Serial)[0].Parallel
	if parallel == nil || len(*parallel) != 2 {
		t.Fatalf("lost iterations: %+v", out.Run)
	}
	for _, call := range *parallel {
		if call.If != nil {
			t.Fatal("parallel iteration inherited another iteration's shadow")
		}
	}
	if (*out.Run.Serial)[1].If == nil {
		t.Fatal("folded predicate using stale outer input")
	}
}

func TestLocalPathsAreNotKnownInputReferences(t *testing.T) {
	for _, path := range []string{"./version", "../version", "/version"} {
		var spec model.OpSpec
		data, err := json.Marshal(map[string]any{"run": map[string]any{"if": []any{map[string]any{"eq": []any{"$(" + path + ")", "old"}}}, "container": map[string]any{"image": map[string]any{"ref": "alpine"}}}})
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(data, &spec); err != nil {
			t.Fatal(err)
		}
		out, err := Expand(&spec, map[string]any{"version": "old"})
		if err != nil {
			t.Fatal(err)
		}
		if out.Run == nil || out.Run.If == nil {
			t.Fatalf("folded local path %s", path)
		}
	}
}

func TestUnresolvedParallelLoopShadowsOuterValues(t *testing.T) {
	var spec model.OpSpec
	if err := json.Unmarshal([]byte(`{"run":{"serial":[{"parallelLoop":{"range":"$(items)","run":{"op":{"ref":"../child","outputs":{"version":""}}}}},{"if":[{"eq":["$(version)","old"]}],"container":{"image":{"ref":"alpine"}}}]}}`), &spec); err != nil {
		t.Fatal(err)
	}
	out, err := Expand(&spec, map[string]any{"version": "old"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Run.Serial == nil || (*out.Run.Serial)[1].If == nil {
		t.Fatal("folded using an input shadowed by an unresolved loop")
	}
}
