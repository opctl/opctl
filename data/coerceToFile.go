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
	// CoerceToFile attempts to coerce value to a file. If coerced, file will be created at rootPath w/ unique name
	CoerceToFile(
		value *model.Value,
		rootPath string,
	) (string, error)
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
	rootPath string,
) (string, error) {
	switch {
	case nil == value:
		path := filepath.Join(rootPath, c.uniqueString.Construct())

		if err := c.ioUtil.WriteFile(
			path,
			[]byte{},
			0666,
		); nil != err {
			return "", fmt.Errorf("Unable to coerce nil to file; error was %v", err.Error())
		}

		return path, nil
	case nil != value.Dir:
		return "", fmt.Errorf("Unable to coerce dir '%v' to file; incompatible types", *value.Dir)
	case nil != value.File:
		return *value.File, nil
	case nil != value.Number:
		path := filepath.Join(rootPath, c.uniqueString.Construct())

		if err := c.ioUtil.WriteFile(
			path,
			[]byte(strconv.FormatFloat(*value.Number, 'f', -1, 64)),
			0666,
		); nil != err {
			return "", fmt.Errorf("Unable to coerce number to file; error was %v", err.Error())
		}

		return path, nil
	case nil != value.Object:
		path := filepath.Join(rootPath, c.uniqueString.Construct())
		valueBytes, err := c.json.Marshal(value.Object)
		if nil != err {
			return "", fmt.Errorf("Unable to coerce object to file; error was %v", err.Error())
		}

		if err := c.ioUtil.WriteFile(
			path,
			valueBytes,
			0666,
		); nil != err {
			return "", fmt.Errorf("Unable to coerce object to file; error was %v", err.Error())
		}

		return path, nil
	case nil != value.String:
		path := filepath.Join(rootPath, c.uniqueString.Construct())

		if err := c.ioUtil.WriteFile(
			path,
			[]byte(*value.String),
			0666,
		); nil != err {
			return "", fmt.Errorf("Unable to coerce string to file; error was %v", err.Error())
		}

		return path, nil
	default:
		return "", fmt.Errorf("Unable to coerce '%#v' to file", value)
	}
}
