package str

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/constraints"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

type Validator interface {
	Validate(
		value *model.Value,
		constraints map[string]interface{},
	) []error
}

func NewValidator() Validator {
	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", constraints.DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", constraints.IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", constraints.SemVerFormatChecker{})

	return _validator{
		coerce: coerce.New(),
	}
}

type _validator struct {
	coerce coerce.Coerce
}

// Validate validates a value against a number parameter
func (vdt _validator) Validate(
	value *model.Value,
	constraints map[string]interface{},
) []error {
	valueAsString, err := vdt.coerce.ToString(value)
	if nil != err {
		return []error{err}
	}

	// guard no constraints
	if nil != constraints {
		errs := []error{}

		valueJsonBytes, err := json.Marshal(valueAsString.String)
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
