package validate

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
)

type Validate interface {
  Param(
  arg *model.Arg,
  param *model.Param,
  ) (errors []error)
}

func New() Validate {
  return &validate{}
}

type validate struct {
}
