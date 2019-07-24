package value

//go:generate counterfeiter -o ./fakeConstructor.go --fake-name FakeConstructor ./ Constructor

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/types"
)

type Constructor interface {
	// Construct constructs a types.Value from an interface{}
	Construct(
		data interface{},
	) (*types.Value, error)
}

func NewConstructor() Constructor {
	return _constructor{}
}

type _constructor struct {
}

func (cvr _constructor) Construct(
	data interface{},
) (*types.Value, error) {
	switch data := data.(type) {
	case bool:
		return &types.Value{Boolean: &data}, nil
	case float64:
		return &types.Value{Number: &data}, nil
	case int:
		// reprocess as float64
		return cvr.Construct(float64(data))
	case string:
		return &types.Value{String: &data}, nil
	case map[string]interface{}:
		return &types.Value{Object: &data}, nil
	case []interface{}:
		return &types.Value{Array: &data}, nil
	default:
		return nil, fmt.Errorf("unable to construct value; '%+v' unexpected type", data)
	}
}
