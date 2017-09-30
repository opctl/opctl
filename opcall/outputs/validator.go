package outputs

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ validator

import (
	"bytes"
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
)

type validator interface {
	Validate(
		outputs map[string]*model.Value,
		params map[string]*model.Param,
	) error
}

func newValidator() validator {
	return _validator{
		data: data.New(),
	}
}

type _validator struct {
	data data.Data
}

func (vdr _validator) Validate(
	outputs map[string]*model.Value,
	params map[string]*model.Param,
) error {
	errMap := map[string][]error{}
	for paramName, paramValue := range params {

		if nil == outputs[paramName] {
			// skip validation of outputs not referenced
			continue
		}

		if errs := vdr.data.Validate(outputs[paramName], paramValue); len(errs) > 0 {
			errMap[paramName] = errs
		}
	}

	if len(errMap) > 0 {
		messageBuffer := bytes.NewBufferString("")
		for outputName, errs := range errMap {
			for _, err := range errs {
				messageBuffer.WriteString(fmt.Sprintf(`
    - %v: %v`, outputName, err.Error()))
			}
		}
		messageBuffer.WriteString(`
`)
		return fmt.Errorf(`
-
  output(s) invalid:
%v
-`, messageBuffer.String())
	}

	return nil
}
