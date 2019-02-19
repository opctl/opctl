package param

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opctl/sdk-golang/model"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

// validateArray validates a value against an array parameter
func (vdt _validator) validateArray(
	value *model.Value,
	constraints map[string]interface{},
) []error {
	valueAsArray, err := vdt.coerce.ToArray(value)
	if nil != err {
		return []error{err}
	}

	// guard no constraints
	if nil != constraints {
		errs := []error{}

		valueJsonBytes, err := json.Marshal(valueAsArray.Array)
		if err != nil {
			// handle syntax errors specially
			return append(
				errs,
				fmt.Errorf("error validating parameter. Details: %v", err.Error()),
			)
		}

		result, err := gojsonschema.Validate(
			gojsonschema.NewGoLoader(constraints),
			gojsonschema.NewBytesLoader(valueJsonBytes),
		)
		if err != nil {
			// handle syntax errors specially
			return append(
				errs,
				fmt.Errorf("error validating param. Details: %v", err.Error()),
			)
		}

		for _, errString := range result.Errors() {
			// enum validation errors include `(root) ` prefix we don't want
			errs = append(errs, errors.New(strings.TrimPrefix(errString.Description(), "(root) ")))
		}

		return errs
	}

	return []error{}
}
