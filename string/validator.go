package string

//go:generate counterfeiter -o ./fakeValidator.go --fake-name fakeValidator ./ Validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chrisdostert/gojsonschema"
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

type Validator interface {
	// Validate validates a string against constraints
	Validate(
		value *string,
		constraints *model.StringConstraints,
	) (errors []error)
}

func newValidator() Validator {
	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", SemVerFormatChecker{})

	return _validator{}
}

type _validator struct {
	os ios.IOS
}

func (this _validator) Validate(
	value *string,
	constraints *model.StringConstraints,
) []error {
	if nil == value {
		return []error{errors.New("string required")}
	}

	// guard no constraints
	if nil != constraints {
		errs := []error{}

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
				fmt.Errorf("Error validating string. Details: %v", err.Error()),
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
				fmt.Errorf("Error validating string. Details: %v", err.Error()),
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
