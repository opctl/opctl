package containerruntime

import "context"

type bootstrapKey struct{}

// WithBootstrap applies the private bootstrap runtime restrictions.
func WithBootstrap(ctx context.Context) context.Context {
	return context.WithValue(WithSeparateOutput(ctx), bootstrapKey{}, true)
}

func BootstrapFrom(ctx context.Context) bool {
	bootstrap, _ := ctx.Value(bootstrapKey{}).(bool)
	return bootstrap
}
