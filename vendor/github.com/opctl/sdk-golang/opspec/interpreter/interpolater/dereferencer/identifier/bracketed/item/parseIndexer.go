package item

//go:generate counterfeiter -o ./fakeParseIndexer.go --fake-name fakeParseIndexer ./ parseIndexer

import (
	"fmt"
	"strconv"
)

// parseIndexer parses identifier as an index of array. If identifier is a negative integer, indexing will occur from the end of the array
type parseIndexer interface {
	ParseIndex(
		identifier string,
		array []interface{},
	) (int64, error)
}

func newParseIndexer() parseIndexer {
	return _parseIndexer{}
}

type _parseIndexer struct {
}

func (dr _parseIndexer) ParseIndex(
	identifier string,
	array []interface{},
) (int64, error) {

	arrayItemIndex, err := strconv.ParseInt(
		identifier,
		10,
		64,
	)
	if nil != err {
		return -1, err
	}

	arrayLength := len(array)
	switch {
	case arrayItemIndex < 0:
		arrayItemIndex = int64(arrayLength) + arrayItemIndex
		if arrayItemIndex < 0 {
			return -1, fmt.Errorf("array index %v out of range 0-%v", arrayItemIndex, arrayLength-1)
		}
	case arrayItemIndex >= int64(arrayLength):
		return -1, fmt.Errorf("array index %v out of range 0-%v", arrayItemIndex, arrayLength-1)
	}

	return arrayItemIndex, nil
}
