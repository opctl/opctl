package containerruntime

import "context"

type separateOutputKey struct{}

// WithSeparateOutput requests distinct stdout and stderr streams for bootstrap.
func WithSeparateOutput(ctx context.Context) context.Context {
	return context.WithValue(ctx, separateOutputKey{}, true)
}

// SeparateOutputFrom reports whether bootstrap needs distinct output streams.
func SeparateOutputFrom(ctx context.Context) bool {
	separate, _ := ctx.Value(separateOutputKey{}).(bool)
	return separate
}
