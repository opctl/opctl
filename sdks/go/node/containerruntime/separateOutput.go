package containerruntime

import "context"

type separateOutputKey struct{}

// WithSeparateOutput requests distinct stdout and stderr process streams without stored runtime logs.
func WithSeparateOutput(ctx context.Context) context.Context {
	return context.WithValue(ctx, separateOutputKey{}, true)
}

// SeparateOutputFrom reports whether distinct process output streams are requested.
func SeparateOutputFrom(ctx context.Context) bool {
	separate, _ := ctx.Value(separateOutputKey{}).(bool)
	return separate
}
