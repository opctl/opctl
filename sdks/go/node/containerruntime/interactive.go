package containerruntime

import (
	"context"

	"github.com/opctl/opctl/sdks/go/model"
)

type interactiveKey struct{}

// WithInteractive sets whether a root op requested interactive containers.
func WithInteractive(ctx context.Context, interactive bool) context.Context {
	return context.WithValue(ctx, interactiveKey{}, interactive)
}

// IsInteractive reports whether the root op requested interactive containers.
func IsInteractive(ctx context.Context) bool {
	interactive, _ := ctx.Value(interactiveKey{}).(bool)
	return interactive
}

type termSizeKey struct{}

// WithInitialTermSize sets the initial console size inherited by container calls.
func WithInitialTermSize(ctx context.Context, size model.TermSize) context.Context {
	return context.WithValue(ctx, termSizeKey{}, size)
}

// InitialTermSize returns the requested size, or zero if the client omitted it.
func InitialTermSize(ctx context.Context) model.TermSize {
	size, _ := ctx.Value(termSizeKey{}).(model.TermSize)
	return size
}

// InteractiveRuntime is implemented by runtimes that can forward terminal input and output.
type InteractiveRuntime interface {
	SupportsInteractive() bool
}
