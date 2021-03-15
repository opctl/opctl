package params

import (
	"bytes"
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
)

//Coerce args to satisfy params
func Coerce(
	values map[string]*model.Value,
	params map[string]*model.Param,
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
		case nil != paramValue.Array:
			coercedValues[paramName], err = coerce.ToArray(value)
		case nil != paramValue.Boolean:
			coercedValues[paramName], err = coerce.ToBoolean(value)
		case nil != paramValue.Dir:
			coercedValues[paramName] = value
		case nil != paramValue.File:
			coercedValues[paramName], err = coerce.ToFile(value, opScratchDir)
		case nil != paramValue.String:
			coercedValues[paramName], err = coerce.ToString(value)
		case nil != paramValue.Number:
			coercedValues[paramName], err = coerce.ToNumber(value)
		case nil != paramValue.Object:
			coercedValues[paramName], err = coerce.ToObject(value)
		case nil != paramValue.Socket:
			coercedValues[paramName] = value
			continue paramLoop
		default:
			err = fmt.Errorf("unable to coerce arg: param was unexpected type %+v", paramValue)
		}

		if nil != err {
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
