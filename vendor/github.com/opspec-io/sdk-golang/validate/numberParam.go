package validate

import (
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/xeipuuv/gojsonschema"
	"math"
	"strings"
)

// validates an value against a string parameter
func (this validate) numberParam(
	rawValue *model.Data,
	param *model.NumberParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue {
		errs = append(errs, errors.New("Number required"))
		return
	}

	value := rawValue.Number
	if 0 == value && 0 != param.Default {
		// apply default if value not set
		value = param.Default
	}

	// guard no constraints
	if paramConstraints := param.Constraints; nil != param.Constraints {

		// perform validations not supported by gojsonschema
		if integerConstraint := paramConstraints.Format; integerConstraint == "integer" {
			if ceiledValue := math.Ceil(value); ceiledValue != value {
				errs = append(errs, fmt.Errorf("Does not match format '%v'", integerConstraint))
			}
		}

		// perform validations supported by gojsonschema
		constraintsJsonBytes, err := format.NewJsonFormat().From(paramConstraints)
		if err != nil {
			// handle syntax errors specially
			errs = append(errs, fmt.Errorf("Error interpreting constraints; the pkg likely has a syntax error. Details: %v", err.Error()))
			return
		}

		valueJsonBytes, err := format.NewJsonFormat().From(value)
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
