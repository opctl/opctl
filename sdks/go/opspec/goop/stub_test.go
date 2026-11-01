package goop

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
)

func TestUntypedArrayStubRoundTrip(t *testing.T) {
	use := StubUse{Name: "items", Ref: "../items"}
	render := func(constraints model.Constraints) []byte {
		data, err := RenderStub(use, &model.OpSpec{Inputs: map[string]*model.ParamSpec{"items": {Array: &model.ArrayParamSpec{Constraints: constraints}}}})
		if err != nil {
			t.Fatal(err)
		}
		return data
	}
	without := render(nil)
	with := render(model.Constraints{"items": map[string]any{}})
	if !bytes.Contains(with, []byte("Items opctl.Array[any]")) || !bytes.Contains(without, []byte("Items opctl.Array[any]")) {
		t.Fatalf("untyped stub changed: %s\n%s", without, with)
	}
}

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
	for _, text := range []string{"\n\n+build ignore", "\n\n  +build ignore", "\n\ngo:build ignore", "\ufeff", "\x00\x01\r\t", "\xff"} {
		t.Run("untrusted "+text, func(t *testing.T) {
			param := &model.ParamSpec{Description: text, String: &model.StringParamSpec{Default: text, Constraints: model.Constraints{"pattern": text}}}
			spec := &model.OpSpec{Description: text, Inputs: map[string]*model.ParamSpec{text: param}, Outputs: map[string]*model.ParamSpec{"result": {String: &model.StringParamSpec{Description: text}}}}
			data, err := RenderStub(StubUse{Name: "greet", Ref: "repo#v1" + text}, spec)
			if err != nil {
				t.Fatal(err)
			}
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "//go:") || strings.HasPrefix(line, "// +build") {
					t.Fatalf("description became a directive: %s", data)
				}
			}
		})
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

func TestRenderStubPreservesUnusualObjectKeys(t *testing.T) {
	for _, key := range []string{"", "-", "a,b", "a\"b", "a\\b", "a'b", "a\tb", "a\ufeffb"} {
		t.Run(key, func(t *testing.T) {
			object := map[string]any{"type": "object", "properties": map[string]any{key: map[string]any{"type": "string"}}, "additionalProperties": map[string]any{"type": "string"}}
			spec := &model.OpSpec{Inputs: map[string]*model.ParamSpec{
				"value":  {Object: &model.ObjectParamSpec{Constraints: model.Constraints{"properties": object["properties"], "additionalProperties": object["additionalProperties"]}}},
				"nested": {Object: &model.ObjectParamSpec{Constraints: model.Constraints{"properties": map[string]any{"child": object}}}},
			}}
			data, err := RenderStub(StubUse{Name: "greet", Ref: "repo#v1"}, spec)
			if err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(string(data), "Value opctl.Object[map[string]any]") || !strings.Contains(string(data), "Child map[string]any") {
				t.Fatalf("object keys would be changed by JSON tags: %s", data)
			}
		})
	}
}

func TestRenderStubOmitsSecretConstraints(t *testing.T) {
	spec := &model.OpSpec{Inputs: map[string]*model.ParamSpec{
		"token": {String: &model.StringParamSpec{IsSecret: true, Default: "secret-default", Constraints: model.Constraints{
			"enum": []string{"secret-enum"}, "const": "secret-const", "pattern": "secret-pattern",
		}}},
		"public": {String: &model.StringParamSpec{Constraints: model.Constraints{"pattern": "public-pattern"}}},
	}}
	data, err := RenderStub(StubUse{Name: "greet", Ref: "repo#v1"}, spec)
	if err != nil {
		t.Fatal(err)
	}
	for _, secret := range []string{"secret-default", "secret-enum", "secret-const", "secret-pattern"} {
		if bytes.Contains(data, []byte(secret)) {
			t.Fatalf("secret constraint or default leaked: %s", data)
		}
	}
	for _, doc := range []string{"Has a default. Secret.", "public-pattern"} {
		if !bytes.Contains(data, []byte(doc)) {
			t.Fatalf("missing public documentation: %s", data)
		}
	}
}
