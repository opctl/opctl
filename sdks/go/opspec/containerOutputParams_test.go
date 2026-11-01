package opspec

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"strings"
	"testing"
)

func TestContainerOutputDeclarationShape(t *testing.T) {
	for _, declarations := range []map[string]*model.ParamSpec{
		{}, {"bad.key": {String: &model.StringParamSpec{}}}, {"x": nil}, {"x": {Number: &model.NumberParamSpec{Default: 0}}}, {"x": {Dir: &model.DirParamSpec{}}}, {"x": {String: &model.StringParamSpec{}, Number: &model.NumberParamSpec{}}},
	} {
		if err := ValidateContainerOutputParams(declarations, map[string]string{"x": "$(x)"}); err == nil {
			t.Fatal("invalid declaration accepted")
		}
	}
	if err := ValidateContainerOutputParams(map[string]*model.ParamSpec{"unbound": {String: &model.StringParamSpec{}}}, map[string]string{"x": "$(x)"}); err != nil {
		t.Fatal(err)
	}
}

func TestContainerOutputRejectsExplicitNullDefaults(t *testing.T) {
	for _, kind := range []string{"string", "number", "boolean", "array", "object"} {
		t.Run(kind, func(t *testing.T) {
			for _, fallback := range []string{"", `"default":null`} {
				var declarations map[string]*model.ParamSpec
				if err := json.Unmarshal([]byte(fmt.Sprintf(`{"result":{"%s":{%s}}}`, kind, fallback)), &declarations); err != nil {
					t.Fatal(err)
				}
				err := ValidateContainerOutputParams(declarations, map[string]string{"result": "$(result)"})
				if fallback == "" {
					if err != nil {
						t.Fatal(err)
					}
				} else if err == nil || !strings.Contains(err.Error(), "cannot have a default") {
					t.Fatalf("explicit null default accepted: %v", err)
				}
			}
		})
	}
}
