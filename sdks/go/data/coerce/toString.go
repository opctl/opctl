package coerce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
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
			return nil, errors.Wrap(err, "unable to coerce array to string")
		}

		arrayBytes, err := json.Marshal(nativeArray)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce array to string")
		}
		arrayString := string(arrayBytes)
		return &model.Value{String: &arrayString}, nil
	case nil != value.Boolean:
		booleanString := strconv.FormatBool(*value.Boolean)
		return &model.Value{String: &booleanString}, nil
	case nil != value.Dir:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce dir '%v' to string", *value.Dir))
	case nil != value.File:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce file to string")
		}
		fileString := string(fileBytes)
		return &model.Value{String: &fileString}, nil
	case nil != value.Number:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		return &model.Value{String: &numberString}, nil
	case nil != value.Object:
		nativeObject, err := value.Unbox()
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce object to string")
		}

		objectBytes, err := json.Marshal(nativeObject)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce object to string")
		}
		objectString := string(objectBytes)
		return &model.Value{String: &objectString}, nil
	case nil != value.String:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to string", value)
	}
}
