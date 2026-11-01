package containerruntime

import (
	"context"
	"io"
)

type bootstrapInputKey struct{}

// BootstrapInputLimit bounds the in-memory credential archive delivered to tmpfs.
const BootstrapInputLimit = 1024 * 1024

// WithBootstrapInput attaches private, node-generated input without serializing it.
func WithBootstrapInput(ctx context.Context, input io.Reader) context.Context {
	return context.WithValue(WithSeparateOutput(ctx), bootstrapInputKey{}, input)
}

func BootstrapInputFrom(ctx context.Context) io.Reader {
	input, _ := ctx.Value(bootstrapInputKey{}).(io.Reader)
	return input
}
