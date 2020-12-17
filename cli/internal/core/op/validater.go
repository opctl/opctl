package op

import (
	"context"
	"fmt"

	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// Validater exposes the "op validate" sub command
type Validater interface {
	Validate(
		ctx context.Context,
		opRef string,
	) (string, error)
}

// newValidater returns an initialized "op validate" sub command
func newValidater(
	dataResolver dataresolver.DataResolver,
) Validater {
	return _validater{
		dataResolver: dataResolver,
	}
}

type _validater struct {
	dataResolver dataresolver.DataResolver
}

func (ivkr _validater) Validate(
	ctx context.Context,
	opRef string,
) (string, error) {
	opDirHandle, err := ivkr.dataResolver.Resolve(
		opRef,
		nil,
	)
	if nil != err {
		return "", err
	}

	successMessage := fmt.Sprintf("%v is valid", opDirHandle.Ref())
	return successMessage, opspec.Validate(ctx, *opDirHandle.Path())
}
