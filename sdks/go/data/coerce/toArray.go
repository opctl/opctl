package coerce

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"io/ioutil"
)

// ToArray coerces a value to an array value
func ToArray(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return nil, nil
	case nil != value.Array:
		return value, nil
	case nil != value.Boolean:
		return nil, fmt.Errorf("unable to coerce boolean '%v' to array; incompatible types", *value.Boolean)
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to array; incompatible types", *value.Dir)
	case nil != value.File:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to array; error was %v", err.Error())
		}
		valueArray := new([]interface{})
		err = json.Unmarshal([]byte(fileBytes), valueArray)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to array; error was %v", err.Error())
		}
		return &model.Value{Array: valueArray}, nil
	case nil != value.Number:
		return nil, fmt.Errorf("unable to coerce number '%v' to array; incompatible types", *value.Number)
	case nil != value.Socket:
		return nil, fmt.Errorf("unable to coerce socket '%v' to array; incompatible types", *value.Socket)
	case nil != value.String:
		valueArray := new([]interface{})
		err := json.Unmarshal([]byte(*value.String), valueArray)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce string to array; error was %v", err.Error())
		}
		return &model.Value{Array: valueArray}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to array", value)
	}
}
