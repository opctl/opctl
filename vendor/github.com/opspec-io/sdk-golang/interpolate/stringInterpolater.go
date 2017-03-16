package interpolate

//go:generate counterfeiter -o ./fakeStringInterpolater.go --fake-name fakeStringInterpolater ./ stringInterpolater

import (
	"fmt"
	"strings"
)

type stringInterpolater interface {
	// interpolates a template w/ a string according to opspec
	Interpolate(s string, varName string, varValue string) string
}

func newStringInterpolater() stringInterpolater {
	return _stringInterpolater{}
}

type _stringInterpolater struct{}

func (this _stringInterpolater) Interpolate(template string, varName string, varValue string) string {
	return strings.Replace(template, fmt.Sprintf(`$(%v)`, varName), varValue, -1)
}
