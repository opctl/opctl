package params

import (
	"strings"
	"testing"

	"github.com/opctl/opctl/sdks/go/model"
)

func TestExactParameterTypes(t *testing.T) {
	text := "3"
	zero := 0.0
	no := false
	empty := ""
	array := []interface{}{}
	object := map[string]interface{}{}
	cases := []struct {
		name    string
		value   *model.Value
		param   *model.ParamSpec
		failure string
	}{
		{"number", &model.Value{Number: &zero}, &model.ParamSpec{Number: &model.NumberParamSpec{}}, ""},
		{"boolean", &model.Value{Boolean: &no}, &model.ParamSpec{Boolean: &model.BooleanParamSpec{}}, ""},
		{"string", &model.Value{String: &empty}, &model.ParamSpec{String: &model.StringParamSpec{}}, ""},
		{"array", &model.Value{Array: &array}, &model.ParamSpec{Array: &model.ArrayParamSpec{}}, ""},
		{"object", &model.Value{Object: &object}, &model.ParamSpec{Object: &model.ObjectParamSpec{}}, ""},
		{"numeric string", &model.Value{String: &text}, &model.ParamSpec{Number: &model.NumberParamSpec{}}, "expected number, received string"},
		{"number as string", &model.Value{Number: &zero}, &model.ParamSpec{String: &model.StringParamSpec{}}, "expected string, received number"},
		{"missing", nil, &model.ParamSpec{String: &model.StringParamSpec{}}, "required output is missing"},
		{"empty union", &model.Value{}, &model.ParamSpec{String: &model.StringParamSpec{}}, "malformed value"},
		{"multiple kinds", &model.Value{String: &text, Number: &zero}, &model.ParamSpec{Number: &model.NumberParamSpec{}}, "malformed value"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Validate(map[string]*model.Value{"count": tc.value}, map[string]*model.ParamSpec{"count": tc.param}, ValidationOptions{ExactTypes: true})
			if tc.failure == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.failure) || !strings.Contains(err.Error(), "count") {
				t.Fatalf("expected %q, received %v", tc.failure, err)
			}
		})
	}
	if err := Validate(map[string]*model.Value{"count": {String: &text}}, map[string]*model.ParamSpec{"count": {Number: &model.NumberParamSpec{}}}); err != nil {
		t.Fatal("legacy numeric string no longer coerces:", err)
	}
}

func TestExactConstraintsAndPrivateErrors(t *testing.T) {
	array := []interface{}{"3"}
	object := map[string]interface{}{"count": "private-invalid-count"}
	for _, tc := range []struct {
		value *model.Value
		param *model.ParamSpec
	}{
		{&model.Value{Array: &array}, &model.ParamSpec{Array: &model.ArrayParamSpec{Constraints: model.Constraints{"items": map[string]interface{}{"type": "number"}}}}},
		{&model.Value{Object: &object}, &model.ParamSpec{Object: &model.ObjectParamSpec{Constraints: model.Constraints{"properties": map[string]interface{}{"count": map[string]interface{}{"type": "number"}}}}}},
	} {
		for _, private := range []bool{false, true} {
			value := *tc.value
			value.Sensitive = private
			err := Validate(map[string]*model.Value{"metadata": &value}, map[string]*model.ParamSpec{"metadata": tc.param}, ValidationOptions{ExactTypes: true})
			if err == nil {
				t.Fatal("nested constraint accepted wrong kind")
			}
			if private && strings.Contains(err.Error(), "private-invalid-count") {
				t.Fatal("private constraint error contains value bytes")
			}
		}
	}
	for _, value := range []*model.Value{nil, {String: new(string)}} {
		err := Validate(map[string]*model.Value{"secret": value}, map[string]*model.ParamSpec{"secret": {Number: &model.NumberParamSpec{IsSecret: true}}}, ValidationOptions{ExactTypes: true})
		if err == nil || !strings.Contains(err.Error(), "sensitive") {
			t.Fatal("declared sensitive failure was not sanitized")
		}
	}
}

func TestExactNestedConstraintLocation(t *testing.T) {
	object := map[string]interface{}{"metadata": map[string]interface{}{"revision": 3}}
	values := map[string]*model.Value{"result": {Object: &object}}
	declarations := map[string]*model.ParamSpec{"result": {Object: &model.ObjectParamSpec{Constraints: model.Constraints{
		"properties": map[string]interface{}{"metadata": map[string]interface{}{
			"properties": map[string]interface{}{"revision": map[string]interface{}{"type": "string"}},
		}},
	}}}}
	err := Validate(values, declarations, ValidationOptions{ExactTypes: true})
	if err == nil || !strings.Contains(err.Error(), "metadata.revision") {
		t.Fatalf("nested constraint location missing: %v", err)
	}
	legacy := Validate(values, declarations)
	if legacy == nil || strings.Contains(legacy.Error(), "metadata.revision") {
		t.Fatalf("legacy constraint diagnostic changed: %v", legacy)
	}
}
