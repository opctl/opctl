package inputs

import "github.com/opspec-io/sdk-golang/model"

func (_inputs _Inputs) Validate(
	inputs map[string]*model.Data,
	params map[string]*model.Param,
) map[string][]error {
	errs := map[string][]error{}
	for paramName, paramValue := range params {
		errs[paramName] = _inputs.validator.Validate(
			inputs[paramName],
			paramValue,
		)
	}
	return errs
}
