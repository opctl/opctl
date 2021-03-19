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
	case value == nil:
		return nil, nil
	case value.Array != nil:
		return nil, errors.Wrap(errIncompatibleTypes, "unable to coerce array to object")
	case value.Dir != nil:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce dir '%v' to object", *value.Dir))
	case value.File != nil:
		fileBytes, err := ioutil.ReadFile(*value.File)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to object")
		}
		valueMap := &map[string]interface{}{}
		err = json.Unmarshal([]byte(fileBytes), valueMap)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce file to object")
		}
		return &model.Value{Object: valueMap}, nil
	case value.Number != nil:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce number '%v' to object", *value.Number))
	case value.Object != nil:
		return value, nil
	case value.String != nil:
		valueMap := &map[string]interface{}{}
		err := json.Unmarshal([]byte(*value.String), valueMap)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce string to object")
		}
		return &model.Value{Object: valueMap}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to object", value)
	}
}
