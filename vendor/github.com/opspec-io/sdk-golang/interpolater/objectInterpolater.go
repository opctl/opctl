package interpolater

//go:generate counterfeiter -o ./fakeObjectInterpolater.go --fake-name fakeObjectInterpolater ./ objectInterpolater

import (
	"github.com/golang-interfaces/encoding-ijson"
)

type objectInterpolater interface {
	// interpolates a template w/ a number according to opspec
	Interpolate(s string, varName string, varValue map[string]interface{}) string
}

func newObjectInterpolater() objectInterpolater {
	return _objectInterpolater{
		stringInterpolater: newStringInterpolater(),
		json:               ijson.New(),
	}
}

type _objectInterpolater struct {
	stringInterpolater stringInterpolater
	json               ijson.IJSON
}

func (this _objectInterpolater) Interpolate(template string, varName string, varValue map[string]interface{}) string {
	// @TODO refactor to support returning errs to caller
	varValueBytes, _ := this.json.Marshal(varValue)
	return this.stringInterpolater.Interpolate(template, varName, string(varValueBytes))
}
