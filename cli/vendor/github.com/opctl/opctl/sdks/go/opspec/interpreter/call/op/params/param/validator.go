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

//counterfeiter:generate -o fakes/validator.go . Validator
type Validator interface {
	Validate(
		value *model.Value,
		param *model.Param,
	) []error
}

func NewValidator() Validator {
	return _validator{
		arrayValidator:   array.NewValidator(),
		booleanValidator: boolean.NewValidator(),
		dirValidator:     dir.NewValidator(),
		fileValidator:    file.NewValidator(),
		numberValidator:  number.NewValidator(),
		objectValidator:  object.NewValidator(),
		strValidator:     str.NewValidator(),
		socketValidator:  socket.NewValidator(),
	}
}

type _validator struct {
	arrayValidator   array.Validator
	booleanValidator boolean.Validator
	dirValidator     dir.Validator
	fileValidator    file.Validator
	numberValidator  number.Validator
	objectValidator  object.Validator
	strValidator     str.Validator
	socketValidator  socket.Validator
}

// Validate validates a value against a parameter
// note: param defaults aren't considered
func (vdt _validator) Validate(
	value *model.Value,
	param *model.Param,
) []error {
	if nil == value {
		return []error{errors.New("required")}
	}

	switch {
	case nil != param.Array:
		return vdt.arrayValidator.Validate(value, param.Array.Constraints)
	case nil != param.Boolean:
		return vdt.booleanValidator.Validate(value)
	case nil != param.Dir:
		return vdt.dirValidator.Validate(value)
	case nil != param.File:
		return vdt.fileValidator.Validate(value)
	case nil != param.Number:
		return vdt.numberValidator.Validate(value, param.Number.Constraints)
	case nil != param.String:
		return vdt.strValidator.Validate(value, param.String.Constraints)
	case nil != param.Object:
		return vdt.objectValidator.Validate(value, param.Object.Constraints)
	case nil != param.Socket:
		return vdt.socketValidator.Validate(value)
	default:
		return []error{fmt.Errorf("unable to validate value; param was unexpected type %+v", param)}
	}
}
