package item

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/value"
)

// Interpreter interprets an item from data via indexString.
// data MUST be an array & indexString MUST parse to a +- integer within bounds of array
//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		indexString string,
		data model.Value,
	) (*model.Value, error)
}

func NewInterpreter() Interpreter {
	return _interpreter{
		parseIndexer:     newParseIndexer(),
		valueConstructor: value.NewConstructor(),
	}
}

type _interpreter struct {
	parseIndexer     parseIndexer
	valueConstructor value.Constructor
}

func (dr _interpreter) Interpret(
	indexString string,
	data model.Value,
) (*model.Value, error) {
	itemIndex, err := dr.parseIndexer.ParseIndex(indexString, *data.Array)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret item; error was %v", err.Error())
	}

	item := (*data.Array)[itemIndex]
	itemValue, err := dr.valueConstructor.Construct(item)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret item; error was %v", err.Error())
	}

	return itemValue, nil
}
