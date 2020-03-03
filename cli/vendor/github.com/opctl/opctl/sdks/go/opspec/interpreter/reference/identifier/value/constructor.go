package value

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o fakes/constructor.go . Constructor
type Constructor interface {
	// Construct constructs a model.Value from an interface{}
	Construct(
		data interface{},
	) (*model.Value, error)
}

func NewConstructor() Constructor {
	return _constructor{}
}

type _constructor struct {
}

func (cvr _constructor) Construct(
	data interface{},
) (*model.Value, error) {
	switch data := data.(type) {
	case bool:
		return &model.Value{Boolean: &data}, nil
	case float64:
		return &model.Value{Number: &data}, nil
	case int:
		// reprocess as float64
		return cvr.Construct(float64(data))
	case string:
		return &model.Value{String: &data}, nil
	case map[string]interface{}:
		return &model.Value{Object: &data}, nil
	case []interface{}:
		return &model.Value{Array: &data}, nil
	default:
		return nil, fmt.Errorf("unable to construct value; '%+v' unexpected type", data)
	}
}
