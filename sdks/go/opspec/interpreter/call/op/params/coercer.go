package params

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeCoercer.go --fake-name FakeCoercer ./ Coercer

import (
	"bytes"
	"fmt"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
)

type Coercer interface {
	// Coerce coerces values for/against params
	Coerce(
		values map[string]*model.Value,
		params map[string]*model.Param,
		opScratchDir string,
	) (
		map[string]*model.Value,
		error,
	)
}

// NewCoercer returns an initialized Coercer instance
func NewCoercer() Coercer {
	return _coercer{
		coerce: coerce.New(),
	}
}

type _coercer struct {
	coerce coerce.Coerce
}

func (crc _coercer) Coerce(
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
		value := values[paramName]
		var err error

		switch {
		case nil != paramValue.Array:
			coercedValues[paramName], err = crc.coerce.ToArray(value)
		case nil != paramValue.Boolean:
			coercedValues[paramName], err = crc.coerce.ToBoolean(value)
		case nil != paramValue.Dir:
			coercedValues[paramName] = value
		case nil != paramValue.File:
			coercedValues[paramName], err = crc.coerce.ToFile(value, opScratchDir)
		case nil != paramValue.String:
			coercedValues[paramName], err = crc.coerce.ToString(value)
		case nil != paramValue.Number:
			coercedValues[paramName], err = crc.coerce.ToNumber(value)
		case nil != paramValue.Object:
			coercedValues[paramName], err = crc.coerce.ToObject(value)
		case nil != paramValue.Socket:
			coercedValues[paramName] = value
			continue paramLoop
		default:
			err = fmt.Errorf("unable to coerce arg; param was unexpected type %+v", paramValue)
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
