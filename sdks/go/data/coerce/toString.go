package coerce

import (
	"bytes"
	stdjson "encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/json"
	"github.com/ipld/go-ipld-prime/node/basicnode"
)

// ToString coerces a value to a string value
func ToString(
	value ipld.Node,
) (ipld.Node, error) {
	nb := basicnode.Prototype.String.NewBuilder()

	switch {
	case value.IsNull():
		nb.AssignString("")
		return nb.Build(), nil
	case value.Kind() == ipld.Kind_List:
		var buf bytes.Buffer
		if err := json.Encode(value, &buf); err != nil {
			return nil, fmt.Errorf("unable to coerce array to string: %w", err)
		}
		nb.AssignString(string(buf.Bytes()))
		return nb.Build(), nil
	case value.Kind() == ipld.Kind_Bool:
		b, err := value.AsBool()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce boolean to string: %w", err)
		}
		nb.AssignString(strconv.FormatBool(b))
		return nb.Build(), nil
	case value.Dir != nil:
		return nil, fmt.Errorf("unable to coerce dir '%v' to string: %w", *value.Dir, errIncompatibleTypes)
	case value.File != nil:
		fileBytes, err := os.ReadFile(*value.File)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce file to string: %w", err)
		}
		fileString := string(fileBytes)
		return ipld.Node{String: &fileString}, nil
	case value.Number != nil:
		numberString := strconv.FormatFloat(*value.Number, 'f', -1, 64)
		return ipld.Node{String: &numberString}, nil
	case value.Object != nil:
		nativeObject, err := value.Unbox()
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to string: %w", err)
		}

		objectBytes, err := json.Marshal(nativeObject)
		if err != nil {
			return nil, fmt.Errorf("unable to coerce object to string: %w", err)
		}
		objectString := string(objectBytes)
		return ipld.Node{String: &objectString}, nil
	case value.String != nil:
		return value, nil
	default:
		return nil, fmt.Errorf("unable to coerce '%+v' to string", value)
	}
}
