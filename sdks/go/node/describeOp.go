package node

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
)

// OpDeclarationError reports invalid declarations or bootstrap program failures.
type OpDeclarationError struct{ Err error }

func (e *OpDeclarationError) Error() string { return e.Err.Error() }
func (e *OpDeclarationError) Unwrap() error { return e.Err }

func (c core) DescribeOp(ctx context.Context, req model.DescribeOpReq) (*model.OpSpec, error) {
	dir, err := c.ResolveData(ctx, req.Op.Ref, req.Op.PullCreds)
	if err != nil {
		return nil, err
	}
	spec, err := opfile.Get(ctx, dir)
	if errors.Is(err, opfile.ErrGoOp) {
		path := dir.Ref()
		if !filepath.IsAbs(path) {
			path = filepath.Join(c.dataCachePath, path)
		}
		return c.bootstrapper.Describe(ctx, &model.OpCall{BaseCall: model.BaseCall{OpPath: path}}, "")
	}
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, &OpDeclarationError{Err: err}
	}
	return spec, nil
}
