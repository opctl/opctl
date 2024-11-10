package coerce

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipld/go-ipld-prime/codec/json"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// ToArray coerces a value to an array value
func ToArray(
	value ipld.Node,
) (ipld.Node, error) {
	nb := basicnode.Prototype.List.NewBuilder()
	
	switch {
	case value.IsNull():
		return nil, errors.New("unable to coerce null to array")
	case value.Kind() == ipld.Kind_List:
		return value, nil
	case value.Kind() == ipld.Kind_Bool:
		return nil, fmt.Errorf("unable to coerce boolean to array: %w", errIncompatibleTypes)
	case value.Kind() == ipld.Kind_String:
		str, err := value.AsString()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce string to array: %w", err)
		}

		strIsPath := filepath.IsAbs(str)
		if strIsPath {
			stat, err := os.Stat(str)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce string to array: %w", err)
			}

			if stat.IsDir() {
				return nil, fmt.Errorf("unable to coerce dir to array: %w", errIncompatibleTypes)
			}

			if !stat.Mode().IsRegular() {
				return nil, fmt.Errorf("unable to coerce socket to array: %w", errIncompatibleTypes)
			}

			file, err := os.Open(str)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce file to array: %w", err)
			}
			defer file.Close()

			if err := json.Decode(nb, file); err != nil {
				return nil, fmt.Errorf("unable to coerce file to array: %w", err)
			}

			return nb.Build(), nil
		}

		if err := json.Decode(nb, strings.NewReader(str)); err != nil {
			return nil, fmt.Errorf("unable to coerce string to array: %w", err)
		}

		return nb.Build(), nil
	case value.Kind() == ipld.Kind_Int || value.Kind() == ipld.Kind_Float:
		return nil, fmt.Errorf("unable to coerce number to array: %w", errIncompatibleTypes)
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to array", value)
	}
}
