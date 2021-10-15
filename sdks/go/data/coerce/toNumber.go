package coerce

import (
	"fmt"
	"os"
	"strconv"

	"github.com/opctl/opctl/sdks/go/model"
)

// ToNumber coerces a value to a number value
func ToNumber(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case value == nil:
		return &model.Value{Number: new(float64)}, nil
	case value.Array != nil:
		return nil, fmt.Errorf("unable to coerce array to number: %w", errIncompatibleTypes)
	case value.Dir != nil:
		return nil, fmt.Errorf("unable to coerce dir to number: %w", errIncompatibleTypes)
	case value.File != nil:
		fileBytes, err := os.ReadFile(*value.File)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to number: %w", err)
		}

		float64Value, err := strconv.ParseFloat(string(fileBytes), 64)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to number: %w", err)
		}
		return &model.Value{Number: &float64Value}, nil
	case value.Number != nil:
		return value, nil
	case value.Object != nil:
		return nil, fmt.Errorf("unable to coerce object to number: %w", errIncompatibleTypes)
	case value.String != nil:
		float64Value, err := strconv.ParseFloat(*value.String, 64)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce string to number: %w", err)
		}
		return &model.Value{Number: &float64Value}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to number", value)
	}
}
