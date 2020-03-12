package opspec

import (
	"context"

	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

//counterfeiter:generate -o fakes/validator.go . Validator
type Validator interface {
	// Validate validates an op
	Validate(
		ctx context.Context,
		opPath string,
	) []error
}

// NewValidator returns an initialized Validator instance
func NewValidator() Validator {
	return _validator{
		opFileGetter: opfile.NewGetter(),
	}
}

type _validator struct {
	opFileGetter opfile.Getter
}

func (vdr _validator) Validate(
	ctx context.Context,
	opPath string,
) []error {
	errs := []error{}
	if _, err := vdr.opFileGetter.Get(
		ctx,
		opPath,
	); nil != err {
		errs = append(errs, err)
	}

	return errs
}
