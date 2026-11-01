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
