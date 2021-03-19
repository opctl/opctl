package item

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
	"github.com/pkg/errors"
)

// Interpret an item from data via indexString.
// data MUST be an array & indexString MUST parse to a +- integer within bounds of array
func Interpret(
	indexString string,
	data model.Value,
) (*model.Value, error) {
	itemIndex, err := ParseIndex(indexString, *data.Array)
	if err != nil {
		return nil, errors.Wrap(err, "unable to interpret item")
	}

	item := (*data.Array)[itemIndex]
	itemValue, err := value.Construct(item)
	if err != nil {
		return nil, errors.Wrap(err, "unable to interpret item")
	}

	return itemValue, nil
}
