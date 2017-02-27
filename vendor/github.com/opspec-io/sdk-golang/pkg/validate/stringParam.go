package validate

import (
	"errors"
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"github.com/opspec-io/sdk-golang/util/format"
	"github.com/xeipuuv/gojsonschema"
	"strings"
)

// validates an value against a string parameter
func (this validate) stringParam(
	rawValue *model.Data,
	param *model.StringParam,
) (errs []error) {
	errs = []error{}

	// handle no value passed
	if nil == rawValue {
		errs = append(errs, errors.New("String required"))
		return
	}

	value := rawValue.String
	if "" == value && "" != param.Default {
		// apply default if value not set
		value = param.Default
	}

	// guard no constraints
	if paramConstraints := param.Constraints; nil != paramConstraints {

		// perform validations supported by gojsonschema
		constraintsJsonBytes, err := format.NewJsonFormat().From(paramConstraints)
		if err != nil {
			// handle syntax errors specially
			errs = append(errs, fmt.Errorf("Error interpreting constraints; the op likely has a syntax error.\n Details: %v", err.Error()))
			return
		}

		valueJsonBytes, err := format.NewJsonFormat().From(value)
		if err != nil {
			// handle syntax errors specially
			errs = append(errs, fmt.Errorf("Error validating parameter.\n Details: %v", err.Error()))
			return
		}

		result, err := gojsonschema.Validate(
			gojsonschema.NewBytesLoader(constraintsJsonBytes),
			gojsonschema.NewBytesLoader(valueJsonBytes),
		)
		if err != nil {
			// handle syntax errors specially
			errs = append(errs, fmt.Errorf("Error interpreting constraints; the op likely has a syntax error.\n Details: %v", err.Error()))
			return
		}

		for _, errString := range result.Errors() {
			// enum validation errors include `(root) ` prefix we don't want
			errs = append(errs, errors.New(strings.TrimPrefix(errString.Description(), "(root) ")))
		}

	}

	return
}
