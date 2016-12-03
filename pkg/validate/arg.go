package validate

import (
  "regexp"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "fmt"
  "unicode/utf8"
)

func (this validate) Arg(
arg *model.Arg,
param *model.Param,
) (errs []error) {
  if ("" != arg.String) {
    errs = this.stringArg(arg.String, param.String)
  }
  return
}

func (this validate) stringArg(
arg string,
param *model.StringParam,
) (errs []error) {
  errs = []error{}
  length := utf8.RuneCountInString(arg)
  if (param.MinLength > 0 && length < param.MinLength) {
    errs = append(errs, fmt.Errorf(
      "%v must be >= %v characters",
      param.Name,
      param.MinLength,
    ))
  }
  if (param.MaxLength > 0 && length > param.MaxLength) {
    errs = append(errs, fmt.Errorf(
      "%v must be <= %v characters",
      param.Name,
      param.MaxLength,
    ))
  }
  if ("" != param.Pattern) {
    isMatch, _ := regexp.MatchString(param.Pattern, arg)
    if (!isMatch) {
      errs = append(errs, fmt.Errorf(
        "%v must match pattern %v",
        param.Name,
        param.Pattern,
      ))
    }
  }
  return
}
