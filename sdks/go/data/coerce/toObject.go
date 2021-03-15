package coerce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

// ToObject coerces a value to an object value
func ToObject(
	value *model.Value,
) (*model.Value, error) {
	switch {
	case nil == value:
		return nil, nil
	case nil != value.Array:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce array to object")
	case nil != value.Dir:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce dir '%v' to object", *value.Dir))
	case nil != value.File:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce file to object")
		}
		valueMap := &map[string]interface{}{}
		err = json.Unmarshal([]byte(fileBytes), valueMap)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce file to object")
		}
		return &model.Value{Object: valueMap}, nil
	case nil != value.Number:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce number '%v' to object", *value.Number))
	case nil != value.Object:
		return value, nil
	case nil != value.String:
		valueMap := &map[string]interface{}{}
		err := json.Unmarshal([]byte(*value.String), valueMap)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce string to object")
		}
		return &model.Value{Object: valueMap}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to object", value)
	}
}
