package sdk

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

type testInputs struct {
	Env     String                 `opctl:"env,description='Target environment, dev or prod',default=dev,enum=dev|prod"`
	Count   Number                 `opctl:"count,default=2,minimum=1,maximum=3,format=integer"`
	Token   String                 `opctl:"token,secret,minLength=1"`
	Src     Dir                    `opctl:"src,default=$(./src)"`
	Config  File                   `opctl:"config,default=$(./config)"`
	Prod    Boolean                `opctl:"wantProd,default=false"`
	Names   Array[string]          `opctl:"names,default='[\"a\"]',minItems=1"`
	Options Object[testOptions]    `opctl:"options"`
	Labels  Object[map[string]int] `opctl:"labels"`
}
type testOptions struct {
	Name   string `json:"name"`
	Counts []int  `json:"counts"`
}

func TestTypedInputs(t *testing.T) {
	data, err := Describe[testInputs]()
	if err != nil {
		t.Fatal(err)
	}
	var response struct {
		Opspec struct{ Inputs map[string]map[string]any }
	}
	if err := json.Unmarshal(data, &response); err != nil {
		t.Fatal(err)
	}
	for _, part := range []string{`"description":"Target environment, dev or prod"`, `"items":{"type":"string"}`, `"additionalProperties":{"type":"integer"}`, `"name":{"type":"string"}`, `"counts":{"items":{"type":"integer"},"type":"array"}`} {
		if !strings.Contains(string(data), part) {
			t.Fatalf("describe missing %s: %s", part, data)
		}
	}
	values := map[string]any{"env": "prod", "count": 3, "wantProd": true, "names": []string{"a", "b"}, "options": map[string]any{"name": "test", "counts": []int{1}}, "labels": map[string]int{"a": 2}}
	data, err = Plan(func(in testInputs) Call {
		if in.Env.Value() != "prod" || in.Count.Value() != 3 || !in.Prod.Value() || in.Names.Value()[1] != "b" || in.Options.Value().Counts[0] != 1 || in.Labels.Value()["a"] != 2 {
			t.Fatal("typed values lost")
		}
		return Container("alpine:3").EnvFrom("TOKEN", in.Token).Dir("/src", in.Src).File("/config", in.Config)
	}, values)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"TOKEN":"$(token)"`) {
		t.Fatal(string(data))
	}
	for _, value := range []string{"token", "src", "config"} {
		bad := map[string]any{}
		for k, v := range values {
			bad[k] = v
		}
		bad[value] = "private"
		if data, err := Plan(func(in testInputs) Call { return Container("alpine") }, bad); err == nil || len(data) != 0 {
			t.Fatalf("accepted forbidden %s", value)
		}
	}
	_, err = Plan(func(in testInputs) Call { return Container(in.Token.Value()) }, values)
	if err == nil || !strings.Contains(err.Error(), `input "token" is secret`) {
		t.Fatalf("secret read: %v", err)
	}
}
func TestInvalidInputDeclarations(t *testing.T) {
	cases := []struct {
		name     string
		describe func() ([]byte, error)
	}{
		{"not struct", Describe[string]},
		{"pointer handle", Describe[struct {
			X *String `opctl:"x"`
		}]},
		{"wrapped handle", Describe[struct {
			X wrappedString `opctl:"x"`
		}]},
		{"no tag", Describe[struct{ X String }]},
		{"unexported", Describe[struct {
			x String `opctl:"x"`
		}]},
		{"plain field", Describe[struct {
			X string `opctl:"x"`
		}]},
		{"unknown", Describe[struct {
			X Number `opctl:"x,minimun=1"`
		}]},
		{"bad name", Describe[struct {
			X String `opctl:"bad.name"`
		}]},
		{"duplicate", Describe[struct {
			X String `opctl:"x"`
			Y String `opctl:"x"`
		}]},
		{"bad key type", Describe[struct {
			X Number `opctl:"x,minLength=1"`
		}]},
		{"bad default", Describe[struct {
			X Number `opctl:"x,default=no"`
		}]},
		{"bad boolean", Describe[struct {
			X Boolean `opctl:"x,default=no"`
		}]},
		{"secret boolean", Describe[struct {
			X Boolean `opctl:"x,secret"`
		}]},
		{"socket default", Describe[struct {
			X Socket `opctl:"x,default=/tmp"`
		}]},
		{"duplicate constraint", Describe[struct {
			X Number `opctl:"x,minimum=1,constraints='{\"minimum\":2}'"`
		}]},
		{"unsupported generic", Describe[struct {
			X Object[int] `opctl:"x"`
		}]},
		{"recursive", Describe[struct {
			X Object[recursive] `opctl:"x"`
		}]},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.describe()
			if err == nil {
				t.Fatal("invalid declaration accepted")
			}
			if tc.name == "duplicate constraint" && !strings.Contains(err.Error(), "duplicate constraint minimum") {
				t.Fatal(err)
			}
		})
	}
}

type recursive struct{ Children []recursive }

func TestTagGrammar(t *testing.T) {
	for _, tc := range []struct {
		tag string
		ok  bool
	}{
		{`x,description='a,b',pattern='a\\b\'c'`, true},
		{"x,description=bare", true}, {"x,description='unterminated", false}, {"x,description='a'junk", false}, {"x,description='bad\\q'", false}, {"x,secret,secret", false}, {"", false},
	} {
		_, err := parseTag(tc.tag)
		if (err == nil) != tc.ok {
			t.Errorf("tag %q: %v", tc.tag, err)
		}
	}
}
func TestBindingsAndLiteralEscaping(t *testing.T) {
	type inputs struct {
		X String `opctl:"x"`
		S Socket `opctl:"s"`
	}
	data, err := Plan(func(in inputs) Call {
		return Container("$(x)").Cmd("$(x)", `\$(x)`, `\\$(x)`).Env("X", in.X.Value()).WorkDir("$(x)").Socket("/socket", in.S)
	}, map[string]any{"x": "$(secret)"})
	if err != nil {
		t.Fatal(err)
	}
	var result response
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatal(err)
	}
	ctr := result.Opspec.Run.Container
	if ctr.Image.Ref != `\$(x)` || ctr.Cmd[1] != `\\\$(x)` || ctr.Cmd[2] != `\\\\\$(x)` || ctr.EnvVars["X"] != `\$(secret)` || ctr.Sockets["/socket"] != "$(s)" {
		t.Fatalf("bad escaping: %s", data)
	}
	for _, call := range []Call{Container("alpine").Env("$X", "value"), Container("alpine").Dir("/$(x)", Dir{}), Container("alpine").EnvFrom("X", String{})} {
		if _, err := Plan(func(struct{}) Call { return call }, nil); err == nil {
			t.Fatal("invalid binding accepted")
		}
	}
}

func TestConstraintTags(t *testing.T) {
	for kind, tag := range map[string]string{
		"string":  `x,description='A,b',default=ok,enum=ok|yes,format=email,pattern='.{1,3}',minLength=1,maxLength=3`,
		"number":  `x,default=2,enum=1|2,format=integer,minimum=1,maximum=3,multipleOf=1`,
		"boolean": `x,default=false`,
		"array":   `x,default='[1,2]',minItems=1,maxItems=3,uniqueItems=true,constraints='{"additionalItems":false}'`,
		"object":  `x,default='{"a":1}',minProperties=1,maxProperties=3,required=a,constraints='{"properties":{"a":{"type":"number"}}}'`,
		"dir":     `x,default=$(./src),secret`, "file": `x,default=$(./config),secret`, "socket": `x,secret`,
	} {
		tags, err := parseTag(tag)
		if err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
		delete(tags, "name")
		var typ reflect.Type
		switch kind {
		case "array":
			typ = reflect.TypeOf([]any{})
		case "object":
			typ = reflect.TypeOf(map[string]any{})
		}
		if _, err := param(kind, typ, tags); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
}

type wrappedString struct{ String }

func TestConstraintsCannotReplaceDerivedSchema(t *testing.T) {
	for _, tc := range []struct {
		key      string
		describe func() ([]byte, error)
	}{
		{"items", Describe[struct {
			X Array[string] `opctl:"x,constraints='{\"items\":{\"type\":\"number\"}}'"`
		}]},
		{"properties", Describe[struct {
			X Object[testOptions] `opctl:"x,constraints='{\"properties\":{}}'"`
		}]},
		{"additionalProperties", Describe[struct {
			X Object[map[string]int] `opctl:"x,constraints='{\"additionalProperties\":false}'"`
		}]},
	} {
		t.Run(tc.key, func(t *testing.T) {
			_, err := tc.describe()
			if err == nil || !strings.Contains(err.Error(), "duplicate constraint "+tc.key) {
				t.Fatalf("expected duplicate constraint %s, got %v", tc.key, err)
			}
		})
	}
}
