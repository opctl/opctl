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
	case value == nil:
		data = []byte{}
	case value.Array != nil:
		nativeArray, err := value.Unbox()
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce array to file")
		}

		data, err = json.Marshal(nativeArray)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce array to file")
		}
	case value.Boolean != nil:
		data = []byte(strconv.FormatBool(*value.Boolean))
	case value.Dir != nil:
		return nil, errors.Wrap(errIncompatibleTypes, fmt.Sprintf("unable to coerce dir '%v' to file", *value.Dir))
	case value.File != nil:
		return value, nil
	case value.Number != nil:
		data = []byte(strconv.FormatFloat(*value.Number, 'f', -1, 64))
	case value.Object != nil:
		nativeObject, err := value.Unbox()
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce object to file")
		}

		data, err = json.Marshal(nativeObject)
		if err != nil {
			return nil, errors.Wrap(err, "unable to coerce object to file")
		}
	case value.String != nil:
		data = []byte(*value.String)
	default:
		data, _ := json.Marshal(value)
		return nil, fmt.Errorf("unable to coerce '%s' to file", string(data))
	}

	uniqueStr, err := uniquestring.Construct()
	if err != nil {
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
		if err = os.MkdirAll(filepath.Dir(path), os.FileMode(0777)); err == nil {
			err = ioutil.WriteFile(
				path,
				data,
				0666,
			)
		}
	}

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("unable to coerce '%+v' to file", value))
	}

	return &model.Value{File: &path}, nil
}
