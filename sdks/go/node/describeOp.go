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
	return describeOpDir(ctx, dir, req.Op.PullCreds, "", "", c.bootstrapper, c.dataCachePath)
}

func describeOpDir(ctx context.Context, dir model.DataHandle, creds *model.Creds, parent, root string, bootstrap bootstrapper, dataCachePath string) (*model.OpSpec, error) {
	spec, err := opfile.Get(ctx, dir)
	if errors.Is(err, opfile.ErrGoOp) {
		path := dir.Ref()
		if !filepath.IsAbs(path) {
			path = filepath.Join(dataCachePath, path)
		}
		op := &model.OpCall{BaseCall: model.BaseCall{OpPath: path}, OpID: parent}
		if fetched, ok := dir.(interface{ RepositoryRoot() string }); ok {
			op.RepositoryRoot = fetched.RepositoryRoot()
		}
		if source, ok := dir.(interface{ SourceResource() string }); ok {
			op.BootstrapSource, op.BootstrapCreds = source.SourceResource(), creds
		}
		return bootstrap.Describe(ctx, op, root)
	}
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, &OpDeclarationError{Err: err}
	}
	return spec, nil
}
