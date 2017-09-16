package data

import (
	"errors"
	"fmt"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

type coerceToNumber interface {
	// CoerceToNumber attempts to coerce value to a number
	CoerceToNumber(
		value *model.Value,
	) (*model.Value, error)
}

func newCoerceToNumber() coerceToNumber {
	return _coerceToNumber{
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _coerceToNumber struct {
	ioUtil iioutil.IIOUtil
	os     ios.IOS
}

func (c _coerceToNumber) CoerceToNumber(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return &model.Value{Number: new(float64)}, nil
	case nil != value.Dir:
		return nil, errors.New("unable to coerce dir to number; incompatible types")
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to number; error was %v", err.Error())
		}

		float64Value, err := strconv.ParseFloat(string(fileBytes), 64)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to number; error was %v", err.Error())
		}
		return &model.Value{Number: &float64Value}, nil
	case nil != value.Number:
		return value, nil
	case nil != value.Object:
		return nil, errors.New("unable to coerce object to number; incompatible types")
	case nil != value.String:
		float64Value, err := strconv.ParseFloat(*value.String, 64)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce string to number; error was %v", err.Error())
		}
		return &model.Value{Number: &float64Value}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to number", value)
	}
}
