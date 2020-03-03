package object

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/constraints"
	"github.com/xeipuuv/gojsonschema"
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

// Validate validates a value against an object parameter
func (vdt _validator) Validate(
	value *model.Value,
	constraints map[string]interface{},
) []error {
	valueAsObject, err := vdt.coerce.ToObject(value)
	if nil != err {
		return []error{err}
	}

	// guard no constraints
	if nil != constraints {
		errs := []error{}

		valueJSONBytes, err := json.Marshal(valueAsObject.Object)
		if err != nil {
			// handle syntax errors specially
			return append(
				errs,
				fmt.Errorf("error validating parameter. Details: %v", err.Error()),
			)
		}

		result, err := gojsonschema.Validate(
			gojsonschema.NewGoLoader(constraints),
			gojsonschema.NewBytesLoader(valueJSONBytes),
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
