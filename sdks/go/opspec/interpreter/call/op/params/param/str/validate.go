package str

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

// Validate validates a value against a number parameter
func Validate(
	value *model.Value,
	constraints map[string]interface{},
) []error {

	// register custom format checkers
	gojsonschema.FormatCheckers.Add("docker-image-ref", formats.DockerImageRefFormatChecker{})
	gojsonschema.FormatCheckers.Add("integer", formats.IntegerFormatChecker{})
	gojsonschema.FormatCheckers.Add("semver", formats.SemVerFormatChecker{})

	valueAsString, err := coerce.ToString(value)
	if err != nil {
		return []error{err}
	}

	// guard no constraints
	if constraints != nil {
		errs := []error{}

		valueJSONBytes, err := json.Marshal(valueAsString.String)
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
			// enum validation errors include `(root) ` prefix we don't want
			errs = append(
				errs,
				errors.New(strings.TrimPrefix(errString.Description(), "(root) ")),
			)
		}

		return errs
	}

	return []error{}
}
