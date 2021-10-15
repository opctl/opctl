package coerce

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/opctl/opctl/sdks/go/model"
)

// ToArray coerces a value to an array value
func ToArray(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case value == nil:
		return nil, errors.New("unable to coerce null to array")
	case value.Array != nil:
		return value, nil
	case value.Boolean != nil:
		return nil, fmt.Errorf("unable to coerce boolean to array: %w", errIncompatibleTypes)
	case value.Dir != nil:
		return nil, fmt.Errorf("unable to coerce dir to array: %w", errIncompatibleTypes)
	case value.File != nil:
		fileBytes, err := os.ReadFile(*value.File)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to array: %w", err)
		}
		valueArray := new([]interface{})
		err = json.Unmarshal([]byte(fileBytes), valueArray)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to array: %w", err)
		}
		return &model.Value{Array: valueArray}, nil
	case value.Number != nil:
		return nil, fmt.Errorf("unable to coerce number to array: %w", errIncompatibleTypes)
	case value.Socket != nil:
		return nil, fmt.Errorf("unable to coerce socket to array: %w", errIncompatibleTypes)
	case value.String != nil:
		valueArray := new([]interface{})
		err := json.Unmarshal([]byte(*value.String), valueArray)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce string to array: %w", err)
		}
		return &model.Value{Array: valueArray}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to array", value)
	}
}
