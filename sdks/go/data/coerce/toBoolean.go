package coerce

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// isStringTruthy ensures value isn't:
// - all "0"'s
// - ""
// - "FALSE" (case insensitive)
// - "F" (case insensitive)
func isStringTruthy(value string) bool {
	stringValueWithoutZeros := strings.Replace(value, "0", "", -1)
	upperCaseStringValue := strings.ToUpper(value)

	return len(stringValueWithoutZeros) != 0 &&
		upperCaseStringValue != "FALSE" &&
		upperCaseStringValue != "F"
}

// ToBoolean attempts to coerce value to a boolean
func ToBoolean(
	value ipld.Node,
) (ipld.Node, error) {
	nb := basicnode.Prototype.Bool.NewBuilder()

	switch {
	case value.IsNull():
		nb.AssignBool(false)
		return nb.Build(), nil
	case value.Kind() == ipld.Kind_List:
		nb.AssignBool(value.Length() != 0)
		return nb.Build(), nil
	case value.Kind() == ipld.Kind_Bool:
		return value, nil
	case value.Kind() == ipld.Kind_String:
		str, err := value.AsString()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce string to boolean: %w", err)
		}

		strIsPath := filepath.IsAbs(str)
		if strIsPath {
			stat, err := os.Stat(str)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce string to boolean: %w", err)
			}

			if stat.IsDir() {
				return nil, fmt.Errorf("unable to coerce dir to boolean: %w", errIncompatibleTypes)
			}

			if !stat.Mode().IsRegular() {
				return nil, fmt.Errorf("unable to coerce socket to boolean: %w", errIncompatibleTypes)
			}

			fileBytes, err := os.ReadFile(str)
			if err != nil {
				return nil, fmt.Errorf("unable to coerce file to boolean: %w", err)
			}

			nb.AssignBool(isStringTruthy(string(fileBytes)))

			return nb.Build(), nil
		}

		nb.AssignBool(isStringTruthy(str))

		return nb.Build(), nil
	case value.Kind() == ipld.Kind_Float:
		nb := basicnode.Prototype.Bool.NewBuilder()
		f, err := value.AsFloat()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce float to boolean: %w", err)
		}

		nb.AssignBool(f == 0)
		return nb.Build(), nil
	case value.Kind() == ipld.Kind_Int:
		nb := basicnode.Prototype.Bool.NewBuilder()
		f, err := value.AsInt()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce int to boolean: %w", err)
		}

		nb.AssignBool(f == 0)
		return nb.Build(), nil
	case value.Kind() == ipld.Kind_Map:
		nb.AssignBool(value.Length() != 0)
		return nb.Build(), nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to boolean", value)
	}
}
