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
	case value == nil:
		return &model.Value{String: new(string)}, nil
	case value.Array != nil:
		nativeArray, err := value.Unbox()
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce array to string")
		}

		arrayBytes, err := json.Marshal(nativeArray)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce array to string")
		}
		arrayString := string(arrayBytes)
		return &model.Value{String: &arrayString}, nil
	case value.Boolean != nil:
		booleanString := strconv.FormatBool(*value.Boolean)
		return &model.Value{String: &booleanString}, nil
	case value.Dir != nil:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce dir '%v' to string", *value.Dir))
	case value.File != nil:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to string")
		}
		fileString := string(fileBytes)
		return &model.Value{String: &fileString}, nil
	case value.Number != nil:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		return &model.Value{String: &numberString}, nil
	case value.Object != nil:
		nativeObject, err := value.Unbox()
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce object to string")
		}

		objectBytes, err := json.Marshal(nativeObject)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce object to string")
		}
		objectString := string(objectBytes)
		return &model.Value{String: &objectString}, nil
	case value.String != nil:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to string", value)
	}
}
