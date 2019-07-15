package coerce

import (
	"fmt"
	"github.com/golang-interfaces/encoding-ijson"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/uniquestring"
	"os"
	"path/filepath"
	"strconv"
)

type toFiler interface {
	// ToFile attempts to coerce value to a file.
	// scratchDir/{UUID} will be used as file path if file creation necessary;
	// if scratchDir doesn't exist it will be created
	ToFile(
		value *model.Value,
		scratchDir string,
	) (*model.Value, error)
}

func newToFiler() toFiler {
	return _toFiler{
		json:         ijson.New(),
		ioUtil:       iioutil.New(),
		os:           ios.New(),
		uniqueString: uniquestring.NewUniqueStringFactory(),
	}
}

type _toFiler struct {
	json         ijson.IJSON
	ioUtil       iioutil.IIOUtil
	os           ios.IOS
	uniqueString uniquestring.UniqueStringFactory
}

func (c _toFiler) ToFile(
	value *model.Value,
	scratchDir string,
) (*model.Value, error) {
	var data []byte

	switch {
	case nil == value:
		data = []byte{}
	case nil != value.Array:
		var err error
		data, err = c.json.Marshal(value.Array)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce array to file; error was %v", err.Error())
		}
	case nil != value.Boolean:
		data = []byte(strconv.FormatBool(*value.Boolean))
	case nil != value.Dir:
		return nil, fmt.Errorf("unable to coerce dir '%v' to file; incompatible types", *value.Dir)
	case nil != value.File:
		return value, nil
	case nil != value.Number:
		data = []byte(strconv.FormatFloat(*value.Number, 'f', -1, 64))
	case nil != value.Object:
		var err error
		data, err = c.json.Marshal(value.Object)
		if nil != err {
			return nil, fmt.Errorf("unable to coerce object to file; error was %v", err.Error())
		}
	case nil != value.String:
		data = []byte(*value.String)
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to file", value)
	}

	uniqueString, err := c.uniqueString.Construct()
	if nil != err {
		return nil, fmt.Errorf("unable to coerce '%+v' to file; error was %v", value, err.Error())
	}

	path := filepath.Join(scratchDir, uniqueString)

	err = c.ioUtil.WriteFile(
		path,
		data,
		0666,
	)
	if os.IsNotExist(err) {
		// ensure path exists & re-attempt
		if err = c.os.MkdirAll(filepath.Dir(path), os.FileMode(0777)); nil == err {
			err = c.ioUtil.WriteFile(
				path,
				data,
				0666,
			)
		}
	}

	if nil != err {
		return nil, fmt.Errorf("unable to coerce '%+v' to file; error was %v", value, err.Error())
	}

	return &model.Value{File: &path}, nil
}
