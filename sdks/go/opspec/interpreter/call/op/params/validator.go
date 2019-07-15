package params

//go:generate counterfeiter -o ./fakeValidator.go --fake-name FakeValidator ./ Validator

import (
	"bytes"
	"fmt"

	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op/params/param"
)

type Validator interface {
	// Validate validates values for/against params
	Validate(
		values map[string]*model.Value,
		params map[string]*model.Param,
	) error
}

// NewValidator returns an initialized Validator instance
func NewValidator() Validator {
	return _validator{
		paramValidator: param.NewValidator(),
	}
}

type _validator struct {
	paramValidator param.Validator
}

func (vdr _validator) Validate(
	values map[string]*model.Value,
	params map[string]*model.Param,
) error {

	paramErrMap := map[string][]error{}
	for paramName, paramValue := range params {
		errs := vdr.paramValidator.Validate(values[paramName], paramValue)
		if len(errs) > 0 {
			paramErrMap[paramName] = errs
		}
	}

	if len(paramErrMap) > 0 {
		// return error w/ fancy formatted msg
		messageBuffer := bytes.NewBufferString("")
		for outputName, errs := range paramErrMap {
			for _, err := range errs {
				messageBuffer.WriteString(fmt.Sprintf(`
    - %v: %v`, outputName, err.Error()))
			}
		}
		messageBuffer.WriteString(`
`)
		return fmt.Errorf(`
-
  validation error(s):
%v
-`, messageBuffer.String())
	}

	return nil
}
