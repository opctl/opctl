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
	if nil == value {
		return []error{errors.New("required")}
	}

	switch {
	case nil != param.Array:
		return array.Validate(value, param.Array.Constraints)
	case nil != param.Boolean:
		return boolean.Validate(value)
	case nil != param.Dir:
		return dir.Validate(value)
	case nil != param.File:
		return file.Validate(value)
	case nil != param.Number:
		return number.Validate(value, param.Number.Constraints)
	case nil != param.String:
		return str.Validate(value, param.String.Constraints)
	case nil != param.Object:
		return object.Validate(value, param.Object.Constraints)
	case nil != param.Socket:
		return socket.Validate(value)
	default:
		return []error{fmt.Errorf("unable to validate value; param was unexpected type %+v", param)}
	}
}
