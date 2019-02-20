package coerce

import (
	"fmt"
	"strings"

	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
)

type toBooleaner interface {
	// ToBoolean attempts to coerce value to a boolean
	ToBoolean(
		value *model.Value,
	) (*model.Value, error)
}

func newToBooleaner() toBooleaner {
	return _toBooleaner{
		os:     ios.New(),
		ioUtil: iioutil.New(),
	}
}

type _toBooleaner struct {
	ioUtil iioutil.IIOUtil
	os     ios.IOS
}

// isStringTruthy ensures value isn't:
// - all "0"'s
// - ""
// - "FALSE" (case insensitive)
// - "F" (case insensitive)
func (c _toBooleaner) isStringTruthy(value string) bool {
	stringValueWithoutZeros := strings.Replace(value, "0", "", -1)
	upperCaseStringValue := strings.ToUpper(value)

	return len(stringValueWithoutZeros) != 0 &&
		upperCaseStringValue != "FALSE" &&
		upperCaseStringValue != "F"
}

func (c _toBooleaner) ToBoolean(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return &model.Value{Boolean: new(bool)}, nil
	case nil != value.Array:
		booleanValue := len(value.Array) != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.Boolean:
		return value, nil
	case nil != value.Dir:
		fileInfos, err := c.ioUtil.ReadDir(*value.Dir)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce dir to boolean; error was %v", err.Error())
		}

		booleanValue := len(fileInfos) > 0
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.File:
		fileBytes, err := c.ioUtil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to boolean; error was %v", err.Error())
		}

		booleanValue := c.isStringTruthy(string(fileBytes))
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.Number:
		booleanValue := *value.Number != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.Object:
		booleanValue := len(value.Object) != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.String:
		booleanValue := c.isStringTruthy(*value.String)
		return &model.Value{Boolean: &booleanValue}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to boolean", value)
	}
}
