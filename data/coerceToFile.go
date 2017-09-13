package data

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"path/filepath"
	"strconv"
)

type coerceToFile interface {
	// CoerceToFile attempts to coerce value to a file.
	// scratchDir/{UUID} will be used as file path if file creation necessary
	CoerceToFile(
		value *model.Value,
		scratchDir string,
	) (*model.Value, error)
}

func newCoerceToFile() coerceToFile {
	return _coerceToFile{
		json:         ijson.New(),
		ioUtil:       iioutil.New(),
		uniqueString: uniquestring.NewUniqueStringFactory(),
	}
}

type _coerceToFile struct {
	json         ijson.IJSON
	ioUtil       iioutil.IIOUtil
	uniqueString uniquestring.UniqueStringFactory
}

func (c _coerceToFile) CoerceToFile(
	value *model.Value,
	scratchDir string,
) (*model.Value, error) {
	switch {
	case nil == value:
		path := filepath.Join(scratchDir, c.uniqueString.Construct())

		if err := c.ioUtil.WriteFile(
			path,
			[]byte{},
			0666,
		); nil != err {
			return nil, fmt.Errorf("unable to coerce nil to file; error was %v", err.Error())
		}

		return &model.Value{File: &path}, nil
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to file; incompatible types", *value.Dir)
	case nil != value.File:
		return value, nil
	case nil != value.Number:
		path := filepath.Join(scratchDir, c.uniqueString.Construct())

		if err := c.ioUtil.WriteFile(
			path,
			[]byte(strconv.FormatFloat(*value.Number, 'f', -1, 64)),
			0666,
		); nil != err {
			return nil, fmt.Errorf("unable to coerce number to file; error was %v", err.Error())
		}

		return &model.Value{File: &path}, nil
	case nil != value.Object:
		path := filepath.Join(scratchDir, c.uniqueString.Construct())
		valueBytes, err := c.json.Marshal(value.Object)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to file; error was %v", err.Error())
		}

		if err := c.ioUtil.WriteFile(
			path,
			valueBytes,
			0666,
		); nil != err {
			return nil, fmt.Errorf("unable to coerce object to file; error was %v", err.Error())
		}

		return &model.Value{File: &path}, nil
	case nil != value.String:
		path := filepath.Join(scratchDir, c.uniqueString.Construct())

		if err := c.ioUtil.WriteFile(
			path,
			[]byte(*value.String),
			0666,
		); nil != err {
			return nil, fmt.Errorf("unable to coerce string to file; error was %v", err.Error())
		}

		return &model.Value{File: &path}, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%#v' to file", value)
	}
}
