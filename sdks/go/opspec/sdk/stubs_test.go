package sdk_test

import (
	"strings"
	"testing"

	opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
	"github.com/opctl/opctl/sdks/go/opspec/sdk/internal/stubgolden/outputs"
)

func TestGeneratedStubPlans(t *testing.T) {
	data, err := opctl.Plan(func(struct{}) opctl.Call {
		c, out := outputs.Call(outputs.Inputs{Text: opctl.StringOf("$(literal)")})
		return opctl.Serial(c, opctl.Container("alpine:3").EnvFrom("TEXT", out.Echoed))
	}, nil)
	if err != nil || !strings.Contains(string(data), `\\$(literal)`) {
		t.Fatalf("%s %v", data, err)
	}
	for _, parallel := range []bool{false, true} {
		_, err := opctl.Plan(func(struct{}) opctl.Call {
			c, out := outputs.Call(outputs.Inputs{Text: opctl.StringOf("hello")})
			check := opctl.Container("alpine:3").EnvFrom("TEXT", out.Echoed)
			if parallel {
				return opctl.Parallel(c, check)
			}
			return opctl.Serial(check, c)
		}, nil)
		if err == nil {
			t.Fatal("accepted invisible stub output")
		}
	}
	_, err = opctl.Plan(func(struct{}) opctl.Call {
		c, _ := outputs.Call(outputs.Inputs{Text: opctl.StringOf("hello")})
		return opctl.Serial(c, c)
	}, nil)
	if err == nil {
		t.Fatal("accepted repeated stub producer")
	}
}
