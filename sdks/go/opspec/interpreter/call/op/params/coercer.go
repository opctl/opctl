package params

import (
	"bytes"
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
)

// Coerce args to satisfy params
func Coerce(
	values map[string]*model.Value,
	params map[string]*model.ParamSpec,
	opScratchDir string,
) (
	map[string]*model.Value,
	error,
) {
	coercedValues := map[string]*model.Value{}

	paramErrMap := map[string]error{}
paramLoop:
	for paramName, paramValue := range params {
		value, ok := values[paramName]
		if !ok {
			// only coerce provided values
			continue
		}

		var err error
		switch {
		case paramValue.Array != nil:
			coercedValues[paramName], err = coerce.ToArray(value)
		case paramValue.Boolean != nil:
			coercedValues[paramName], err = coerce.ToBoolean(value)
		case paramValue.Dir != nil:
			coercedValues[paramName] = value
		case paramValue.File != nil:
			coercedValues[paramName], err = coerce.ToFile(value, opScratchDir)
		case paramValue.String != nil:
			coercedValues[paramName], err = coerce.ToString(value)
		case paramValue.Number != nil:
			coercedValues[paramName], err = coerce.ToNumber(value)
		case paramValue.Object != nil:
			coercedValues[paramName], err = coerce.ToObject(value)
		case paramValue.Socket != nil:
			coercedValues[paramName] = value
			continue paramLoop
		default:
			err = fmt.Errorf("unable to coerce arg: param was unexpected type %+v", paramValue)
		}

		if err != nil {
			paramErrMap[paramName] = err
		}

	}

	if len(paramErrMap) > 0 {
		// return error w/ fancy formatted msg
		messageBuffer := bytes.NewBufferString("")
		for outputName, err := range paramErrMap {
			messageBuffer.WriteString(fmt.Sprintf(`
    	- %v: %v`, outputName, err))
		}
		messageBuffer.WriteString(`
`)
		return coercedValues, fmt.Errorf(`
-
  validation error(s):
%v
-`, messageBuffer.String())
	}

	return coercedValues, nil
}
