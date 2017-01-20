package validate

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"regexp"
	"unicode/utf8"
)

// validates an arg against a parameter
func (this validate) Param(
	arg *model.Data,
	param *model.Param,
) (errs []error) {
	if nil == param {
		panic("param required")
	}

	switch {
	case nil != param.String:
		errs = this.stringParam(arg, param.String)
	case nil != param.Socket:
		errs = this.socketParam(arg, param.Socket)
	}
	return
}

// validates an arg against a string parameter
func (this validate) stringParam(
	rawArg *model.Data,
	param *model.StringParam,
) (errs []error) {
	errs = []error{}

	// handle no arg passed
	if nil == rawArg {
		errs = append(errs, fmt.Errorf("%v required", param.Name))
		return
	}

	arg := rawArg.String
	if "" == arg && "" != param.Default {
		// apply default if arg not set
		arg = param.Default
	}

	// guard no constraints
	if nil == param.Constraints {
		return
	}
	lengthConstraint := param.Constraints.Length
	if nil != lengthConstraint {
		length := utf8.RuneCountInString(arg)
		if lengthConstraint.Min > 0 && length < lengthConstraint.Min {
			errs = append(errs, fmt.Errorf(
				"%v must be >= %v characters",
				param.Name,
				lengthConstraint.Min,
			))
		}
		if lengthConstraint.Max > 0 && length > lengthConstraint.Max {
			errs = append(errs, fmt.Errorf(
				"%v must be <= %v characters",
				param.Name,
				lengthConstraint.Max,
			))
		}
	}
	patternConstraints := param.Constraints.Patterns
	if len(patternConstraints) > 0 {
		for _, patternConstraint := range patternConstraints {
			isMatch, _ := regexp.MatchString(patternConstraint.Regex, arg)
			if !isMatch {
				errs = append(errs, fmt.Errorf(
					"%v must match pattern %v",
					param.Name,
					patternConstraint.Regex,
				))
			}
		}
	}
	return
}

// validates an arg against a network socket parameter
func (this validate) socketParam(
	rawArg *model.Data,
	param *model.SocketParam,
) (errs []error) {
	errs = []error{}

	// handle no arg passed
	if nil == rawArg || "" == rawArg.Socket {
		errs = append(errs, fmt.Errorf("%v required", param.Name))
	}
	return
}
