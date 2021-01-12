package coerce

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/opctl/opctl/sdks/go/model"
)

// ToPort coerces a value to a port value
func ToPort(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return &model.Value{Port: new(uint16)}, nil
	case nil != value.Array:
		return nil, errors.New("unable to coerce array to port; incompatible types")
	case nil != value.Dir:
		return nil, errors.New("unable to coerce dir to port; incompatible types")
	case nil != value.Number:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		uint64Number, err := strconv.ParseUint(numberString, 10, 16)
		uint16Number := uint16(uint64Number)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce number to port; error was %v", err.Error())
		}
		return &model.Value{Port: &uint16Number}, nil
	case nil != value.String:
		uint64Number, err := strconv.ParseUint(*value.String, 10, 16)
		uint16Number := uint16(uint64Number)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce string to port; error was %v", err.Error())
		}
		return &model.Value{Port: &uint16Number}, nil
	case nil != value.Port:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to port", value)
	}
}
