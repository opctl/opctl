package coerce

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
)

type toObjecter interface {
	// ToObject attempts to coerce value to an object
	ToObject(
		value *model.Value,
	) (*model.Value, error)
}

func newToObjecter() toObjecter {
	return _toObjecter{
		json:   ijson.New(),
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _toObjecter struct {
	ioUtil iioutil.IIOUtil
	json   ijson.IJSON
	os     ios.IOS
}

func (c _toObjecter) ToObject(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return nil, nil
	case nil != value.Array:
		return nil, fmt.Errorf("unable to coerce array to object; incompatible types")
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to object; incompatible types", *value.Dir)
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to object; error was %v", err.Error())
		}
		valueMap := map[string]interface{}{}
		err = c.json.Unmarshal([]byte(fileBytes), &valueMap)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to object; error was %v", err.Error())
		}
		return &model.Value{Object: valueMap}, nil
	case nil != value.Number:
		return nil, fmt.Errorf("unable to coerce number '%v' to object; incompatible types", *value.Number)
	case nil != value.Object:
		return value, nil
	case nil != value.String:
		valueMap := map[string]interface{}{}
		err := c.json.Unmarshal([]byte(*value.String), &valueMap)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce string to object; error was %v", err.Error())
		}
		return &model.Value{Object: valueMap}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to object", value)
	}
}
