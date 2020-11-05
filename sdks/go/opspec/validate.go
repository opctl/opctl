package opspec

import (
	"context"

	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// Validate an op
func Validate(
	ctx context.Context,
	opPath string,
) error {
	_, err := opfile.Get(
		ctx,
		opPath,
	)

	return err
}
