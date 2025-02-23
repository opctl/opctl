package jsonschema

import (
	"errors"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param/formats"
	"github.com/xeipuuv/gojsonschema"
)

// Validate validates a value against constraints
func Validate(
	value any,
	schema map[string]interface{},
) []error {
	// guard no schema
	if schema == nil {
		return nil
	}

	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", formats.DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", formats.IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", formats.SemVerFormatChecker{})

	result, err := gojsonschema.Validate(
		gojsonschema.NewGoLoader(schema),
		gojsonschema.NewGoLoader(value),
	)
	if err != nil {
		return []error{
			err,
		}
	}

	errs := []error{}
	for _, verr := range result.Errors() {
		// enum validation errors include `(root) ` prefix we don't want
		errs = append(
			errs,
			errors.New(strings.TrimPrefix(verr.Description(), "(root) ")),
		)
	}

	return errs
}
