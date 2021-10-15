package coerce

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/opctl/opctl/sdks/go/model"
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
			return nil, fmt.Errorf("unable to coerce array to string: %w", err)
		}

		arrayBytes, err := json.Marshal(nativeArray)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce array to string: %w", err)
		}
		arrayString := string(arrayBytes)
		return &model.Value{String: &arrayString}, nil
	case value.Boolean != nil:
		booleanString := strconv.FormatBool(*value.Boolean)
		return &model.Value{String: &booleanString}, nil
	case value.Dir != nil:
		return nil, fmt.Errorf("unable to coerce dir '%v' to string: %w", *value.Dir, errIncompatibleTypes)
	case value.File != nil:
		fileBytes, err := os.ReadFile(*value.File)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to string: %w", err)
		}
		fileString := string(fileBytes)
		return &model.Value{String: &fileString}, nil
	case value.Number != nil:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		return &model.Value{String: &numberString}, nil
	case value.Object != nil:
		nativeObject, err := value.Unbox()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to string: %w", err)
		}

		objectBytes, err := json.Marshal(nativeObject)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to string: %w", err)
		}
		objectString := string(objectBytes)
		return &model.Value{String: &objectString}, nil
	case value.String != nil:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to string", value)
	}
}
