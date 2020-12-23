package value

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
)

// Construct constructs a model.Value from an interface{}
func Construct(
	data interface{},
) (*model.Value, error) {
	switch data := data.(type) {
	case bool:
		return &model.Value{Boolean: &data}, nil
	case float64:
		return &model.Value{Number: &data}, nil
	case int:
		// reprocess as float64
		return Construct(float64(data))
	case string:
		return &model.Value{String: &data}, nil
	case map[string]interface{}:
		return &model.Value{Object: &data}, nil
	case []interface{}:
		return &model.Value{Array: &data}, nil
	default:
		return nil, fmt.Errorf("unable to construct value; '%v' unexpected type", data)
	}
}
