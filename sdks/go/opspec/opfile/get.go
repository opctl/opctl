package opfile

import (
	"context"
	"io"

	"github.com/opctl/opctl/sdks/go/model"
)

// Get gets the validated, deserialized representation of an "op.yml" file
func Get(
	ctx context.Context,
	opDir model.DataHandle,
) (
	*model.OpSpec,
	error,
) {
	rs, err := opDir.GetContent(
		ctx,
		FileName,
	)
	if err != nil {
		return nil, err
	}

	opFileBytes, err := io.ReadAll(rs)
	if err != nil {
		return nil, err
	}

	return Unmarshal(opFileBytes)
}
