package opspec

import (
	"github.com/opctl/opctl/sdks/go/model"
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
