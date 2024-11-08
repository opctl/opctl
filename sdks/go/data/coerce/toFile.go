package coerce

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/internal/unsudo"
	"github.com/opctl/opctl/sdks/go/model"
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
			return nil, fmt.Errorf("unable to coerce array to file: %w", err)
		}

		data, err = json.Marshal(nativeArray)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce array to file: %w", err)
		}
	case value.Boolean != nil:
		data = []byte(strconv.FormatBool(*value.Boolean))
	case value.Dir != nil:
		return nil, fmt.Errorf("unable to coerce dir '%v' to file: %w", *value.Dir, errIncompatibleTypes)
	case value.File != nil:
		return value, nil
	case value.Number != nil:
		data = []byte(strconv.FormatFloat(*value.Number, 'f', -1, 64))
	case value.Object != nil:
		nativeObject, err := value.Unbox()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to file: %w", err)
		}

		data, err = json.Marshal(nativeObject)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to file: %w", err)
		}
	case value.String != nil:
		data = []byte(*value.String)
	default:
		data, _ := json.Marshal(value)
		return nil, fmt.Errorf("unable to coerce '%s' to file", string(data))
	}

	uniqueStr, err := uniquestring.Construct()
	if err != nil {
		return nil, fmt.Errorf("unable to coerce '%+v' to file: %w", value, err)
	}

	path := filepath.Join(scratchDir, uniqueStr)

	if err := unsudo.CreateFile(path, data); err != nil {
		return nil, fmt.Errorf("unable to coerce '%+v' to file: %w", value, err)
	}

	return &model.Value{File: &path}, nil
}
