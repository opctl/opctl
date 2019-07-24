package coerce

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/sdks/go/types"
	"strconv"
)

type toStringer interface {
	// ToString attempts to coerce value to a string
	ToString(
		value *types.Value,
	) (*types.Value, error)
}

func newToStringer() toStringer {
	return _toStringer{
		json:   ijson.New(),
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _toStringer struct {
	ioUtil iioutil.IIOUtil
	json   ijson.IJSON
	os     ios.IOS
}

func (c _toStringer) ToString(
	value *types.Value,
) (*types.Value, error) {
	switch {
	case nil == value:
		return &types.Value{String: new(string)}, nil
	case nil != value.Array:
		arrayBytes, err := c.json.Marshal(value.Array)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce array to string; error was %v", err.Error())
		}
		arrayString := string(arrayBytes)
		return &types.Value{String: &arrayString}, nil
	case nil != value.Boolean:
		booleanString := strconv.FormatBool(*value.Boolean)
		return &types.Value{String: &booleanString}, nil
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to string; incompatible types", *value.Dir)
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to string; error was %v", err.Error())
		}
		fileString := string(fileBytes)
		return &types.Value{String: &fileString}, nil
	case nil != value.Number:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		return &types.Value{String: &numberString}, nil
	case nil != value.Object:
		objectBytes, err := c.json.Marshal(value.Object)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to string; error was %v", err.Error())
		}
		objectString := string(objectBytes)
		return &types.Value{String: &objectString}, nil
	case nil != value.String:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to string", value)
	}
}
