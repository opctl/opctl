package validate

import (
  "regexp"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "fmt"
  "unicode/utf8"
  "errors"
)

// validates the provided parameter
func (this validate) Param(
param *model.Param,
) (errs []error) {
  if (nil != param.String) {
    errs = this.stringParam(param.String)
  } else if (nil != param.NetSocket) {
    errs = this.netSocketParam(param.NetSocket)
  }
  return
}

// validates the provided string parameter
func (this validate) stringParam(
param *model.StringParam,
) (errs []error) {
  errs = []error{}

  value := param.Value
  if ("" == value && "" != param.Default) {
    // apply default if value not set
    value = param.Default
  }

  // guard no constraints
  if (nil == param.Constraints) {
    return
  }
  lengthConstraint := param.Constraints.Length
  if (nil != lengthConstraint) {
    length := utf8.RuneCountInString(value)
    if (lengthConstraint.Min > 0 && length < lengthConstraint.Min) {
      errs = append(errs, fmt.Errorf(
        "%v must be >= %v characters",
        param.Name,
        lengthConstraint.Min,
      ))
    }
    if (lengthConstraint.Max > 0 && length > lengthConstraint.Max) {
      errs = append(errs, fmt.Errorf(
        "%v must be <= %v characters",
        param.Name,
        lengthConstraint.Max,
      ))
    }
  }
  patternConstraints := param.Constraints.Patterns
  if (len(patternConstraints) > 0) {
    for _, patternConstraint := range patternConstraints {
      isMatch, _ := regexp.MatchString(patternConstraint.Regex, value)
      if (!isMatch) {
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

// validates the provided network socket parameter
func (this validate) netSocketParam(
param *model.NetSocketParam,
) (errs []error) {
  errs = []error{}
  value := param.Value
  if ("" == value.Host) {
    errs = append(errs, errors.New("Host required"))
  }
  if (0 >= value.Port) {
    errs = append(errs, errors.New("Port must be > 0"))
  }
  if (65536 <= value.Port) {
    errs = append(errs, errors.New("Port must be <= 65535"))
  }

  // guard no constraints
  if (nil == param.Constraints) {
    return
  }
  portNumberConstraint := param.Constraints.PortNumber
  if ( nil != portNumberConstraint) {
    if (portNumberConstraint.Number != value.Port) {
      errs = append(errs, fmt.Errorf(
        "%v Port must be %v",
        param.Name,
        value.Port,
      ))
    }
  }
  return
}
