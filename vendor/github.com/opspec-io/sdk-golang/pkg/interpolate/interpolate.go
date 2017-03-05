package interpolate

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Interpolate

import "github.com/opspec-io/sdk-golang/pkg/model"

type Interpolate interface {
	// interpolates a template according to opspec
	Interpolate(template string, scope map[string]*model.Data) string
}

func New() Interpolate {
	return interpolate{
		numberInterpolater: newNumberInterpolater(),
		stringInterpolater: newStringInterpolater(),
	}
}

type interpolate struct {
	numberInterpolater numberInterpolater
	stringInterpolater stringInterpolater
}

// O(n) complexity (n being len(scope))
func (this interpolate) Interpolate(
	template string,
	scope map[string]*model.Data,
) string {
	for varName, varData := range scope {
		if nil != varData {
			switch {
			case 0 != varData.Number:
				template = this.numberInterpolater.Interpolate(template, varName, varData.Number)
			case "" != varData.String:
				template = this.stringInterpolater.Interpolate(template, varName, varData.String)
			}
		}
	}
	return template
}
