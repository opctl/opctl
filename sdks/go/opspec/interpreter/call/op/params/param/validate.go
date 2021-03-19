package param

import (
	"errors"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/array"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/boolean"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/dir"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/file"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/number"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/object"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/socket"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/str"
)

// Validate validates a value against a parameter
// note: param defaults aren't considered
func Validate(
	value *model.Value,
	param *model.Param,
) []error {
	if value == nil {
		return []error{errors.New("required")}
	}

	switch {
	case param.Array != nil:
		return array.Validate(value, param.Array.Constraints)
	case param.Boolean != nil:
		return boolean.Validate(value)
	case param.Dir != nil:
		return dir.Validate(value)
	case param.File != nil:
		return file.Validate(value)
	case param.Number != nil:
		return number.Validate(value, param.Number.Constraints)
	case param.String != nil:
		return str.Validate(value, param.String.Constraints)
	case param.Object != nil:
		return object.Validate(value, param.Object.Constraints)
	case param.Socket != nil:
		return socket.Validate(value)
	default:
		return []error{fmt.Errorf("unable to validate value: param was unexpected type %+v", param)}
	}
}
