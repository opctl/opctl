package object

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ Validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chrisdostert/gojsonschema"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

type Validator interface {
	// Validate validates an object against constraints
	Validate(
		value map[string]interface{},
		constraints *model.ObjectConstraints,
	) []error
}

func newValidator() Validator {
	return _validator{}
}

type _validator struct{}

func (this _validator) Validate(
	value map[string]interface{},
	constraints *model.ObjectConstraints,
) []error {
	if nil == value {
		return []error{errors.New("object required")}
	}

	// guard no constraints
	if nil != constraints {
		errs := []error{}

		// perform validations supported by gojsonschema
		constraintsJsonBytes, err := json.Marshal(constraints)
		if err != nil {
			// handle syntax errors specially
			return append(
				errs,
				fmt.Errorf("Error interpreting constraints; the pkg likely has a syntax error. Details: %v", err.Error()),
			)
		}

		valueJsonBytes, err := json.Marshal(value)
		if err != nil {
			// handle syntax errors specially
			return append(
				errs,
				fmt.Errorf("Error validating object. Details: %v", err.Error()),
			)
		}

		result, err := gojsonschema.Validate(
			gojsonschema.NewBytesLoader(constraintsJsonBytes),
			gojsonschema.NewBytesLoader(valueJsonBytes),
		)
		if err != nil {
			// handle syntax errors specially
			return append(
				errs,
				fmt.Errorf("Error validating object. Details: %v", err.Error()),
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
