package container

import (
	"github.com/opctl/opctl/sdks/go/model"
	"testing"
)

func TestProgrammaticContainerOutputDeclarations(t *testing.T) {
	for _, p := range []*model.ParamSpec{nil, {Dir: &model.DirParamSpec{}}, {Number: &model.NumberParamSpec{Default: 0}}, {Number: &model.NumberParamSpec{}, String: &model.StringParamSpec{}}} {
		spec := &model.ContainerCallSpec{Image: &model.ContainerCallImageSpec{Ref: "alpine"}, Outputs: map[string]string{"count": "$(count)"}, OutputParams: map[string]*model.ParamSpec{"count": p}}
		if _, err := Interpret(nil, spec, "id", "/op", t.TempDir()); err == nil {
			t.Fatal("invalid programmatic declaration accepted")
		}
	}
}
