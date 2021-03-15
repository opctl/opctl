package coerce

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

// ToFile attempts to coerce value to a file
func ToFile(
	value *model.Value,
	scratchDir string,
) (*model.Value, error) {
	var data []byte

	switch {
	case nil == value:
		data = []byte{}
	case nil != value.Array:
		nativeArray, err := value.Unbox()
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce array to file")
		}

		data, err = json.Marshal(nativeArray)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce array to file")
		}
	case nil != value.Boolean:
		data = []byte(strconv.FormatBool(*value.Boolean))
	case nil != value.Dir:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce dir '%v' to file", *value.Dir))
	case nil != value.File:
		return value, nil
	case nil != value.Number:
		data = []byte(strconv.FormatFloat(*value.Number, 'f', -1, 64))
	case nil != value.Object:
		nativeObject, err := value.Unbox()
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce object to file")
		}

		data, err = json.Marshal(nativeObject)
		if nil != err {
			return nil, errors.Wrap(err, "unable to coerce object to file")
		}
	case nil != value.String:
		data = []byte(*value.String)
	default:
		data, _ := json.Marshal(value)
		return nil, fmt.Errorf("unable to coerce '%s' to file", string(data))
	}

	uniqueStr, err := uniquestring.Construct()
	if nil != err {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to coerce '%+v' to file", value))
	}

	path := filepath.Join(scratchDir, uniqueStr)

	err = ioutil.WriteFile(
		path,
		data,
		0666,
	)
	if os.IsNotExist(err) {
		// ensure path exists & re-attempt
		if err = os.MkdirAll(filepath.Dir(path), os.FileMode(0777)); nil == err {
			err = ioutil.WriteFile(
				path,
				data,
				0666,
			)
		}
	}

	if nil != err {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to coerce '%+v' to file", value))
	}

	return &model.Value{File: &path}, nil
}
