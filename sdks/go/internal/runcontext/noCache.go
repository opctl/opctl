// Package runcontext carries execution settings shared by calls in one run.
package runcontext

import "context"

type noCacheRootKey struct{}

// WithNoCacheRoot makes a run mount its cache dirs under noCacheRootPath instead of the cache root.
func WithNoCacheRoot(ctx context.Context, noCacheRootPath string) context.Context {
	return context.WithValue(ctx, noCacheRootKey{}, noCacheRootPath)
}

// NoCacheRoot returns where a --no-cache run mounts its cache dirs, or "" to use the cache root.
func NoCacheRoot(ctx context.Context) string {
	noCacheRootPath, _ := ctx.Value(noCacheRootKey{}).(string)
	return noCacheRootPath
}
