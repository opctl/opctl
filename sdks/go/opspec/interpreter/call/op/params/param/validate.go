package param

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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
	param *model.ParamSpec,
	exact ...bool,
) []error {
	if value == nil {
		if len(exact) > 0 && exact[0] {
			return []error{errors.New("required output is missing")}
		}
		return []error{errors.New("required")}
	}

	if len(exact) > 0 && exact[0] {
		actual, expected := populatedKind(value), populatedKind(param)
		if actual == "malformed value" || actual != expected {
			return []error{fmt.Errorf("expected %s, received %s", expected, actual)}
		}
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

func populatedKind(union interface{}) string {
	value := reflect.ValueOf(union)
	if value.IsNil() {
		return "malformed value"
	}
	value = value.Elem()
	kind, count := "", 0
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Kind() == reflect.Pointer && !field.IsNil() {
			count++
			kind = strings.ToLower(value.Type().Field(i).Name)
		}
	}
	if count != 1 {
		return "malformed value"
	}
	return kind
}
