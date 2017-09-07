package inputs

import (
	"errors"
	"github.com/opspec-io/sdk-golang/model"
)

func (_inputs _Inputs) Validate(
	inputs map[string]*model.Value,
	params map[string]*model.Param,
) map[string][]error {
	errMap := map[string][]error{}
	for paramName, paramValue := range params {

		if nil == inputs[paramName] {
			errMap[paramName] = []error{errors.New("param required")}
			continue
		}

		if errs := _inputs.data.Validate(inputs[paramName], paramValue); len(errs) > 0 {
			errMap[paramName] = errs
		}
	}
	return errMap
}
