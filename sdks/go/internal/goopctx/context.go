// Package goopctx marks execution inside a Go op graph without changing YAML-only runs.
package goopctx

import "context"

type key struct{}

func WithActive(ctx context.Context) context.Context { return context.WithValue(ctx, key{}, true) }
func Active(ctx context.Context) bool                { active, _ := ctx.Value(key{}).(bool); return active }
