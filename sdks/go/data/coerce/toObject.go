package coerce

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime"
)

// ToObject coerces a value to an object value
func ToObject(
	value ipld.Node,
) (ipld.Node, error) {
	switch {
	case value == nil:
		return nil, nil
	case value.Array != nil:
		return nil, fmt.Errorf("unable to coerce array to object: %w", errIncompatibleTypes)
	case value.Dir != nil:
		return nil, fmt.Errorf("unable to coerce dir '%v' to object: %w", *value.Dir, errIncompatibleTypes)
	case value.File != nil:
		fileBytes, err := os.ReadFile(*value.File)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to object: %w", err)
		}
		valueMap := &map[string]interface{}{}
		err = json.Unmarshal([]byte(fileBytes), valueMap)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to object: %w", err)
		}
		return ipld.Node{Object: valueMap}, nil
	case value.Number != nil:
		return nil, fmt.Errorf("unable to coerce number '%v' to object: %w", *value.Number, errIncompatibleTypes)
	case value.Object != nil:
		return value, nil
	case value.String != nil:
		valueMap := &map[string]interface{}{}
		err := json.Unmarshal([]byte(*value.String), valueMap)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce string to object: %w", err)
		}
		return ipld.Node{Object: valueMap}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to object", value)
	}
}
