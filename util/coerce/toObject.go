package coerce

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
)

type toObject interface {
	// ToObject attempts to coerce value to an object
	ToObject(
		value *model.Value,
	) (map[string]interface{}, error)
}

func newToObject() toObject {
	return _toObject{
		json:   ijson.New(),
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _toObject struct {
	ioUtil iioutil.IIOUtil
	json   ijson.IJSON
	os     ios.IOS
}

func (c _toObject) ToObject(
	value *model.Value,
) (map[string]interface{}, error) {
	switch {
	case nil == value:
		return nil, nil
	case nil != value.Dir:
		return nil, fmt.Errorf("Unable to coerce dir '%v' to object; incompatible types", *value.Dir)
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("Unable to coerce file to object; error was %v", err.Error())
		}
		valueMap := map[string]interface{}{}
		err = c.json.Unmarshal([]byte(fileBytes), &valueMap)
		if nil != err {
			return nil, fmt.Errorf("Unable to coerce file to object; error was %v", err.Error())
		}
		return valueMap, nil
	case nil != value.Number:
		return nil, fmt.Errorf("Unable to coerce number '%v' to object; incompatible types", *value.Number)
	case nil != value.Object:
		return value.Object, nil
	case nil != value.String:
		valueMap := map[string]interface{}{}
		err := c.json.Unmarshal([]byte(*value.String), &valueMap)
		if nil != err {
			return nil, fmt.Errorf("Unable to coerce string to object; error was %v", err.Error())
		}
		return valueMap, nil
	default:
		return nil, fmt.Errorf("Unable to coerce '%#v' to object", value)
	}
}
