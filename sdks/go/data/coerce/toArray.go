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
	case value == nil:
		return nil, errors.New("unable to coerce null to array")
	case value.Array != nil:
		return value, nil
	case value.Boolean != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce boolean to array")
	case value.Dir != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce dir to array")
	case value.File != nil:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to array")
		}
		valueArray := new([]interface{})
		err = json.Unmarshal([]byte(fileBytes), valueArray)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to array")
		}
		return &model.Value{Array: valueArray}, nil
	case value.Number != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce number to array")
	case value.Socket != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce socket to array")
	case value.String != nil:
		valueArray := new([]interface{})
		err := json.Unmarshal([]byte(*value.String), valueArray)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce string to array")
		}
		return &model.Value{Array: valueArray}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to array", value)
	}
}
