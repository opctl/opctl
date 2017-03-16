package interpolate

//go:generate counterfeiter -o ./fakeNumberInterpolater.go --fake-name fakeNumberInterpolater ./ numberInterpolater

import (
	"strconv"
)

type numberInterpolater interface {
	// interpolates a template w/ a number according to opspec
	Interpolate(s string, varName string, varValue float64) string
}

func newNumberInterpolater() numberInterpolater {
	return _numberInterpolater{
		stringInterpolater: newStringInterpolater(),
	}
}

type _numberInterpolater struct {
	stringInterpolater stringInterpolater
}

func (this _numberInterpolater) Interpolate(template string, varName string, varValue float64) string {
	return this.stringInterpolater.Interpolate(template, varName, strconv.FormatFloat(varValue, 'f', -1, 64))
}
