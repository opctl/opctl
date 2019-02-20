package coerce

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
)

type toArrayer interface {
	// ToArray attempts to coerce value to an array
	ToArray(
		value *model.Value,
	) (*model.Value, error)
}

func newToArrayer() toArrayer {
	return _toArrayer{
		json:   ijson.New(),
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _toArrayer struct {
	ioUtil iioutil.IIOUtil
	json   ijson.IJSON
	os     ios.IOS
}

func (c _toArrayer) ToArray(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return nil, nil
	case nil != value.Array:
		return value, nil
	case nil != value.Boolean:
		return nil, fmt.Errorf("unable to coerce boolean '%v' to array; incompatible types", *value.Boolean)
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to array; incompatible types", *value.Dir)
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to array; error was %v", err.Error())
		}
		valueArray := []interface{}{}
		err = c.json.Unmarshal([]byte(fileBytes), &valueArray)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to array; error was %v", err.Error())
		}
		return &model.Value{Array: valueArray}, nil
	case nil != value.Number:
		return nil, fmt.Errorf("unable to coerce number '%v' to array; incompatible types", *value.Number)
	case nil != value.Socket:
		return nil, fmt.Errorf("unable to coerce socket '%v' to array; incompatible types", *value.Socket)
	case nil != value.String:
		valueArray := []interface{}{}
		err := c.json.Unmarshal([]byte(*value.String), &valueArray)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce string to array; error was %v", err.Error())
		}
		return &model.Value{Array: valueArray}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to array", value)
	}
}
