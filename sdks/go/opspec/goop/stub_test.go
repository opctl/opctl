package goop

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
)

func TestRenderStub(t *testing.T) {
	var spec model.OpSpec
	if err := json.Unmarshal([]byte(`{"description":"Greets someone.\n//go:generate evil\u0000", "inputs": {
	 "name":{"string":{}}, "token":{"string":{"isSecret":true,"default":"do-not-render"}},
	 "db-name":{"number":{}},"db_name":{"boolean":{}},"123":{"file":{}},"!!!":{"socket":{}},"src":{"dir":{}},
	 "tags":{"array":{"constraints":{"items":{"type":"string"}},"default":[]}},
	 "config":{"object":{"constraints":{"properties":{"style":{"type":"string"},"loud":{"type":"boolean"},"nested":{"type":"object","properties":{"count":{"type":"integer"}}}},"required":["style"]},"default":{"style":"plain"}}},
	 "map":{"object":{"constraints":{"additionalProperties":{"type":"number"}}}},
	 "fallback":{"array":{"constraints":{"items":{"oneOf":[{"type":"string"},{"type":"number"}]}}}},
	 "a-b":{"object":{"constraints":{"properties":{"c":{"type":"object","properties":{"x":{"type":"string"}}}}}}},
	 "a":{"object":{"constraints":{"properties":{"b-c":{"type":"object","properties":{"y":{"type":"number"}}}}}}}
	 },"outputs":{"greeting":{"string":{}},"data":{"object":{}},"items":{"array":{}}}}`), &spec); err != nil {
		t.Fatal(err)
	}
	use := StubUse{Name: "greet", Ref: "../greet", ModuleRef: ".opspec/greet", Local: true}
	data, err := RenderStub(use, &spec)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(data, []byte("do-not-render")) || bytes.Contains(data, []byte{0}) || strings.Contains(string(data), "\n//go:") {
		t.Fatalf("unsafe stub: %s", data)
	}
	if !strings.Contains(string(data), "DbName_") || !strings.Contains(string(data), "InputsABC_") || !strings.Contains(string(data), "omitempty") {
		t.Fatalf("missing stable names or optional tags: %s", data)
	}
	path := "../sdk/internal/stubgolden/greet/stub.go"
	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("../sdk/internal/stubgolden/greet", 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(data, want) {
		t.Fatalf("stub differs from golden:\n%s", data)
	}
	for i := 0; i < 5; i++ {
		next, err := RenderStub(use, &spec)
		if err != nil || !bytes.Equal(data, next) {
			t.Fatal("nondeterministic rendering", err)
		}
	}
}

func TestRenderSharedDescribe(t *testing.T) {
	for _, name := range []string{"inputs", "outputs", "metadata", "docgreet", "docmerge"} {
		data, err := os.ReadFile("../sdk/testdata/describe." + name + ".json")
		if err != nil {
			t.Fatal(err)
		}
		spec, err := ParseDescribe(data)
		if err != nil {
			t.Fatal(err)
		}
		data, err = RenderStub(StubUse{Name: name, Ref: "repo#v1"}, spec)
		if err != nil {
			t.Fatal(err)
		}
		path := "../sdk/internal/stubgolden/" + name + "/stub.go"
		if os.Getenv("UPDATE_GOLDEN") == "1" {
			if err := os.MkdirAll("../sdk/internal/stubgolden/"+name, 0755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, data, 0644); err != nil {
				t.Fatal(err)
			}
		}
		want, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, want) {
			t.Fatalf("%s changed", name)
		}
	}
}
