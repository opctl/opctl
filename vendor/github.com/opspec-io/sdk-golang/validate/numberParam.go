package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/xeipuuv/gojsonschema"
	"math"
	"strings"
)

// validates an value against a string parameter
func (this validate) numberParam(
	rawValue *float64,
	param *model.NumberParam,
) (errs []error) {
	errs = []error{}

	value := rawValue
	if nil == value && nil != param.Default {
		// apply default if value not set
		value = param.Default
	}

	if nil == value {
		errs = append(errs, errors.New("Number required"))
		return
	}

	// guard no constraints
	if paramConstraints := param.Constraints; nil != param.Constraints {

		// perform validations not supported by gojsonschema
		if integerConstraint := paramConstraints.Format; integerConstraint == "integer" {
			if ceiledValue := math.Ceil(*value); ceiledValue != *value {
				errs = append(errs, fmt.Errorf("Does not match format '%v'", integerConstraint))
			}
		}

		// perform validations supported by gojsonschema
		constraintsJsonBytes, err := json.Marshal(paramConstraints)
		if err != nil {
			// handle syntax errors specially
			errs = append(errs, fmt.Errorf("Error interpreting constraints; the pkg likely has a syntax error. Details: %v", err.Error()))
			return
		}

		valueJsonBytes, err := json.Marshal(value)
		if err != nil {
			// handle syntax errors specially
			errs = append(errs, fmt.Errorf("Error validating parameter. Details: %v", err.Error()))
			return
		}

		result, err := gojsonschema.Validate(
			gojsonschema.NewBytesLoader(constraintsJsonBytes),
			gojsonschema.NewBytesLoader(valueJsonBytes),
		)
		if err != nil {
			// handle syntax errors specially
			errs = append(errs, fmt.Errorf("Error validating param. Details: %v", err.Error()))
			return
		}

		for _, errString := range result.Errors() {
			// enum validation errors include `(root) ` prefix we don't want
			errs = append(errs, errors.New(strings.TrimPrefix(errString.Description(), "(root) ")))
		}
	}

	return
}
