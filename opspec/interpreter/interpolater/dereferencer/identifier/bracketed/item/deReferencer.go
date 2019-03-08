package item

//go:generate counterfeiter -o ./fakeDeReferencer.go --fake-name FakeDeReferencer ./ DeReferencer

import (
	"fmt"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer/identifier/value"
)

// DeReferencer dereferences an item from data via indexString.
// data MUST be an array & indexString MUST parse to a +- integer within bounds of array
type DeReferencer interface {
	DeReference(
		indexString string,
		data model.Value,
	) (*model.Value, error)
}

func NewDeReferencer() DeReferencer {
	return _deReferencer{
		parseIndexer:     newParseIndexer(),
		valueConstructor: value.NewConstructor(),
	}
}

type _deReferencer struct {
	parseIndexer     parseIndexer
	valueConstructor value.Constructor
}

func (dr _deReferencer) DeReference(
	indexString string,
	data model.Value,
) (*model.Value, error) {
	itemIndex, err := dr.parseIndexer.ParseIndex(indexString, data.Array)
	if nil != err {
		return nil, fmt.Errorf("unable to deReference item; error was %v", err.Error())
	}

	item := data.Array[itemIndex]
	itemValue, err := dr.valueConstructor.Construct(item)
	if nil != err {
		return nil, fmt.Errorf("unable to deReference item; error was %v", err.Error())
	}

	return itemValue, nil
}
