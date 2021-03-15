package coerce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

// ToArray coerces a value to an array value
func ToArray(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return nil, errors.New("unable to coerce null to array")
	case nil != value.Array:
		return value, nil
	case nil != value.Boolean:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce boolean to array")
	case nil != value.Dir:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce dir to array")
	case nil != value.File:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce file to array")
		}
		valueArray := new([]interface{})
		err = json.Unmarshal([]byte(fileBytes), valueArray)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce file to array")
		}
		return &model.Value{Array: valueArray}, nil
	case nil != value.Number:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce number to array")
	case nil != value.Socket:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce socket to array")
	case nil != value.String:
		valueArray := new([]interface{})
		err := json.Unmarshal([]byte(*value.String), valueArray)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce string to array")
		}
		return &model.Value{Array: valueArray}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to array", value)
	}
}
