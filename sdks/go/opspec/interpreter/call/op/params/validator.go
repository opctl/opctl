package params

import (
	"bytes"
	"fmt"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/sensitivity"
	"sort"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op/params/param"
)

// ValidationOptions opts into strict validation without changing legacy coercion.
type ValidationOptions struct{ ExactTypes bool }

// Validate validates values for/against params
func Validate(
	values map[string]*model.Value,
	params map[string]*model.ParamSpec,
	options ...ValidationOptions,
) error {

	exact := len(options) > 0 && options[0].ExactTypes
	if exact {
		names := make([]string, 0, len(params))
		for name := range params {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, name := range names {
			errs := param.Validate(values[name], params[name], true)
			if len(errs) == 0 {
				continue
			}
			if sensitivity.Secret(params[name]) || (values[name] != nil && values[name].Sensitive) {
				return fmt.Errorf("%q: sensitive parameter failed validation", name)
			}
			reasons := make([]string, len(errs))
			for i, err := range errs {
				reasons[i] = err.Error()
			}
			return fmt.Errorf("%q: %s", name, strings.Join(reasons, "; "))
		}
		return nil
	}

	paramErrMap := map[string][]error{}
	for paramName, paramValue := range params {
		errs := param.Validate(values[paramName], paramValue)
		if len(errs) > 0 && values[paramName] != nil && values[paramName].Sensitive {
			errs = []error{fmt.Errorf("sensitive parameter %q failed validation", paramName)}
		}
		if len(errs) > 0 {
			paramErrMap[paramName] = errs
		}
	}

	if len(paramErrMap) > 0 {
		// return error w/ fancy formatted msg
		messageBuffer := bytes.NewBufferString("")
		for paramName, errs := range paramErrMap {
			for _, err := range errs {
				messageBuffer.WriteString(fmt.Sprintf(`
    - %v: %v`, paramName, err.Error()))
			}
		}
		messageBuffer.WriteString(`
`)
		return fmt.Errorf(`
-
  validation error(s):
%v
-`, messageBuffer.String())
	}

	return nil
}
