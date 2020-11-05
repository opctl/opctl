package item

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
)

// Interpret an item from data via indexString.
// data MUST be an array & indexString MUST parse to a +- integer within bounds of array
func Interpret(
	indexString string,
	data model.Value,
) (*model.Value, error) {
	itemIndex, err := ParseIndex(indexString, *data.Array)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret item; error was %v", err.Error())
	}

	item := (*data.Array)[itemIndex]
	itemValue, err := value.Construct(item)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret item; error was %v", err.Error())
	}

	return itemValue, nil
}
