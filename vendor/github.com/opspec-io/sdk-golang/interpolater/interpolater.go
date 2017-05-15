// Package interpolater implements an interpolater for string templates
package interpolater

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Interpolater

import "github.com/opspec-io/sdk-golang/model"

type Interpolater interface {
	// Interpolate interpolates a string template
	Interpolate(template string, scope map[string]*model.Data) string
}

func New() Interpolater {
	return _Interpolater{
		numberInterpolater: newNumberInterpolater(),
		stringInterpolater: newStringInterpolater(),
	}
}

type _Interpolater struct {
	numberInterpolater numberInterpolater
	stringInterpolater stringInterpolater
}

// Interpolate interpolates a string template
// O(n) complexity (n being len(scope))
func (this _Interpolater) Interpolate(
	template string,
	scope map[string]*model.Data,
) string {
	for varName, varData := range scope {
		if nil != varData {
			switch {
			case nil != varData.Number:
				template = this.numberInterpolater.Interpolate(template, varName, *varData.Number)
			case nil != varData.String:
				template = this.stringInterpolater.Interpolate(template, varName, *varData.String)
			}
		}
	}
	return template
}
