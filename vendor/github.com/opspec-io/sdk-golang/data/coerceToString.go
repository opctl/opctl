package data

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"strconv"
)

type coerceToString interface {
	// CoerceToString attempts to coerce value to a string
	CoerceToString(
		value *model.Value,
	) (*model.Value, error)
}

func newCoerceToString() coerceToString {
	return _coerceToString{
		json:   ijson.New(),
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _coerceToString struct {
	ioUtil iioutil.IIOUtil
	json   ijson.IJSON
	os     ios.IOS
}

func (c _coerceToString) CoerceToString(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return &model.Value{String: new(string)}, nil
	case nil != value.Array:
		arrayBytes, err := c.json.Marshal(value.Array)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce array to string; error was %v", err.Error())
		}
		arrayString := string(arrayBytes)
		return &model.Value{String: &arrayString}, nil
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to string; incompatible types", *value.Dir)
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to string; error was %v", err.Error())
		}
		fileString := string(fileBytes)
		return &model.Value{String: &fileString}, nil
	case nil != value.Number:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		return &model.Value{String: &numberString}, nil
	case nil != value.Object:
		objectBytes, err := c.json.Marshal(value.Object)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to string; error was %v", err.Error())
		}
		objectString := string(objectBytes)
		return &model.Value{String: &objectString}, nil
	case nil != value.String:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to string", value)
	}
}
