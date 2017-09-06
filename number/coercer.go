package number

//go:generate counterfeiter -o ./fakeCoercer.go --fake-name fakeCoercer ./ coercer

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

type coercer interface {
	// Coerce attempts to coerce value to a number
	Coerce(
		value *model.Value,
	) (float64, error)
}

func newCoercer() coercer {
	return _coercer{
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _coercer struct {
	ioUtil iioutil.IIOUtil
	os     ios.IOS
}

func (c _coercer) Coerce(
	value *model.Value,
) (float64, error) {
	switch {
	case nil == value:
		return 0, nil
	case nil != value.Dir:
		return 0, errors.New("Unable to coerce dir to number; incompatible types")
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return 0, fmt.Errorf("Unable to coerce file to number; error was %v", err.Error())
		}

		float64Value, err := strconv.ParseFloat(string(fileBytes), 64)
		if nil != err {
			return 0, fmt.Errorf("Unable to coerce file to number; error was %v", err.Error())
		}
		return float64Value, nil
	case nil != value.Number:
		return *value.Number, nil
	case nil != value.Object:
		return 0, errors.New("Unable to coerce object to number; incompatible types")
	case nil != value.String:
		float64Value, err := strconv.ParseFloat(*value.String, 64)
		if nil != err {
			return 0, fmt.Errorf("Unable to coerce string to number; error was %v", err.Error())
		}
		return float64Value, nil
	default:
		return 0, fmt.Errorf("Unable to coerce '%v' to number", value)
	}
}
