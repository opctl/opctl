package coerce

import (
	"fmt"
	"strings"

	"io/ioutil"

	"github.com/opctl/opctl/sdks/go/model"
)

// isStringTruthy ensures value isn't:
// - all "0"'s
// - ""
// - "FALSE" (case insensitive)
// - "F" (case insensitive)
func isStringTruthy(value string) bool {
	stringValueWithoutZeros := strings.Replace(value, "0", "", -1)
	upperCaseStringValue := strings.ToUpper(value)

	return len(stringValueWithoutZeros) != 0 &&
		upperCaseStringValue != "FALSE" &&
		upperCaseStringValue != "F"
}

// ToBoolean attempts to coerce value to a boolean
func ToBoolean(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return &model.Value{Boolean: new(bool)}, nil
	case nil != value.Array:
		booleanValue := nil != value.Array && len(*value.Array) != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.Boolean:
		return value, nil
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir to boolean; incompatible types")
	case nil != value.File:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce file to boolean; error was %v", err.Error())
		}

		booleanValue := isStringTruthy(string(fileBytes))
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.Number:
		booleanValue := *value.Number != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.Object:
		booleanValue := nil != value.Object && len(*value.Object) != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case nil != value.String:
		booleanValue := isStringTruthy(*value.String)
		return &model.Value{Boolean: &booleanValue}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to boolean", value)
	}
}
