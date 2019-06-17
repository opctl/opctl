package op

//go:generate counterfeiter -o ./fakeValidator.go --fake-name FakeValidator ./ Validator

import (
	"context"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/opfile"
)

type Validator interface {
	// Validate validates an op
	Validate(
		ctx context.Context,
		opHandle model.DataHandle,
	) []error
}

// NewValidator returns an initialized Validator instance
func NewValidator() Validator {
	return _validator{
		dotYmlGetter: dotyml.NewGetter(),
	}
}

type _validator struct {
	dotYmlGetter dotyml.Getter
}

func (vdr _validator) Validate(
	ctx context.Context,
	opHandle model.DataHandle,
) []error {
	errs := []error{}
	if _, err := vdr.dotYmlGetter.Get(
		ctx,
		opHandle,
	); nil != err {
		errs = append(errs, err)
	}

	return errs
}
