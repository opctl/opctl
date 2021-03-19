package coerce

import (
	"fmt"
	"strings"

	"io/ioutil"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
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
	case value == nil:
		return &model.Value{Boolean: new(bool)}, nil
	case value.Array != nil:
		booleanValue := value.Array != nil && len(*value.Array) != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case value.Boolean != nil:
		return value, nil
	case value.Dir != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce dir to boolean")
	case value.File != nil:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to boolean")
		}

		booleanValue := isStringTruthy(string(fileBytes))
		return &model.Value{Boolean: &booleanValue}, nil
	case value.Number != nil:
		booleanValue := *value.Number != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case value.Object != nil:
		booleanValue := value.Object != nil && len(*value.Object) != 0
		return &model.Value{Boolean: &booleanValue}, nil
	case value.String != nil:
		booleanValue := isStringTruthy(*value.String)
		return &model.Value{Boolean: &booleanValue}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to boolean", value)
	}
}
