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
	case nil == value:
		return &model.Value{Number: new(float64)}, nil
	case nil != value.Array:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce array to number")
	case nil != value.Dir:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce dir to number")
	case nil != value.File:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce file to number")
		}

		float64Value, err := strconv.ParseFloat(string(fileBytes), 64)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce file to number")
		}
		return &model.Value{Number: &float64Value}, nil
	case nil != value.Number:
		return value, nil
	case nil != value.Object:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce object to number")
	case nil != value.String:
		float64Value, err := strconv.ParseFloat(*value.String, 64)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce string to number")
		}
		return &model.Value{Number: &float64Value}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to number", value)
	}
}
