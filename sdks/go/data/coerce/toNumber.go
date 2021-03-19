package coerce

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

// ToNumber coerces a value to a number value
func ToNumber(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case value == nil:
		return &model.Value{Number: new(float64)}, nil
	case value.Array != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce array to number")
	case value.Dir != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce dir to number")
	case value.File != nil:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to number")
		}

		float64Value, err := strconv.ParseFloat(string(fileBytes), 64)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to number")
		}
		return &model.Value{Number: &float64Value}, nil
	case value.Number != nil:
		return value, nil
	case value.Object != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce object to number")
	case value.String != nil:
		float64Value, err := strconv.ParseFloat(*value.String, 64)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce string to number")
		}
		return &model.Value{Number: &float64Value}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to number", value)
	}
}
