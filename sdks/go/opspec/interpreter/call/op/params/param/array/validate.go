package array

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/formats"
	"github.com/xeipuuv/gojsonschema"
)

// Validate validates a value against a string parameter
func Validate(
	value *model.Value,
	constraints map[string]interface{},
	includeLocation ...bool,
) []error {
	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", formats.DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", formats.IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", formats.SemVerFormatChecker{})

	valueAsArray, err := coerce.ToArray(value)
	if err != nil {
		return []error{err}
	}

	unboxedValueAsArray, err := valueAsArray.Unbox()
	if err != nil {
		return []error{
			err,
		}
	}

	// guard no constraints
	if constraints != nil {
		errs := []error{}

		valueJSONBytes, err := json.Marshal(unboxedValueAsArray)
		if err != nil {
			// handle syntax errors specially
			return append(
				errs,
				fmt.Errorf("error validating parameter: %w", err),
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
				fmt.Errorf("error validating parameter: %w", err),
			)
		}

		for _, errString := range result.Errors() {
			// Keep legacy diagnostics unchanged unless exact output validation requests locations.
			description := strings.TrimPrefix(errString.Description(), "(root) ")
			if len(includeLocation) > 0 && includeLocation[0] && errString.Field() != "(root)" {
				description = errString.Field() + ": " + description
			}
			errs = append(errs, errors.New(description))
		}

		return errs
	}

	return []error{}
}
