package coerce

import (
	"encoding/json"
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"io/ioutil"
	"strconv"
)

// ToString coerces a value to a string value
func ToString(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return &model.Value{String: new(string)}, nil
	case nil != value.Array:
		nativeArray, err := value.Unbox()
		if nil != err {
			return nil, fmt.Errorf("unable to coerce array to string; error was %v", err)
		}

		arrayBytes, err := json.Marshal(nativeArray)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce array to string; error was %v", err.Error())
		}
		arrayString := string(arrayBytes)
		return &model.Value{String: &arrayString}, nil
	case nil != value.Boolean:
		booleanString := strconv.FormatBool(*value.Boolean)
		return &model.Value{String: &booleanString}, nil
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to string; incompatible types", *value.Dir)
	case nil != value.File:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to string; error was %v", err.Error())
		}
		fileString := string(fileBytes)
		return &model.Value{String: &fileString}, nil
	case nil != value.Number:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		return &model.Value{String: &numberString}, nil
	case nil != value.Object:
		nativeObject, err := value.Unbox()
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to string; error was %v", err)
		}

		objectBytes, err := json.Marshal(nativeObject)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to string; error was %v", err.Error())
		}
		objectString := string(objectBytes)
		return &model.Value{String: &objectString}, nil
	case nil != value.String:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to string", value)
	}
}
