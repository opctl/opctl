package opfile

import (
	"context"
	"errors"
	"io"

	"github.com/opctl/opctl/sdks/go/model"
)

// ErrGoOp signals that a Go op needs bootstrap instead of YAML decoding.
var ErrGoOp = errors.New("Go ops are not supported yet")

// ErrBothOpFiles signals an ambiguous op declaration.
var ErrBothOpFiles = errors.New("op dir contains both op.yml and op.go; an op must be a YAML op or a Go op, not both")

// IsGoOp detects Go ops and rejects op dirs containing both declarations.
func IsGoOp(ctx context.Context, opDir model.DataHandle) (bool, error) {
	yaml, yamlErr := opDir.GetContent(ctx, FileName)
	if yamlErr == nil {
		yaml.Close()
	}
	goFile, goErr := opDir.GetContent(ctx, "op.go")
	if goErr != nil {
		return false, nil
	}
	goFile.Close()
	if yamlErr == nil {
		return false, ErrBothOpFiles
	}
	return true, nil
}

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
	goFile, goErr := opDir.GetContent(ctx, "op.go")
	if goErr == nil {
		goFile.Close()
		if err == nil {
			rs.Close()
			return nil, ErrBothOpFiles
		}
		return nil, ErrGoOp
	}
	if err != nil {
		return nil, err
	}
	defer rs.Close()

	opFileBytes, err := io.ReadAll(rs)
	if err != nil {
		return nil, err
	}

	return Unmarshal(opFileBytes)
}
