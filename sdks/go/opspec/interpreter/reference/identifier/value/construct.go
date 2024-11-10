package value

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
)

// Construct constructs a ipld.Node from an interface{}
func Construct(
	data interface{},
) (*ipld.Node, error) {
	switch data := data.(type) {
	case bool:
		return &ipld.Node{Boolean: &data}, nil
	case float64:
		return &ipld.Node{Number: &data}, nil
	case int:
		// reprocess as float64
		return Construct(float64(data))
	case string:
		return &ipld.Node{String: &data}, nil
	case map[string]interface{}:
		return &ipld.Node{Object: &data}, nil
	case []interface{}:
		return &ipld.Node{Array: &data}, nil
	default:
		return nil, fmt.Errorf("unable to construct value: '%v' unexpected type", data)
	}
}
