package node

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/opctl/opctl/sdks/go/model"
)

type privateCompletionKey struct{}

// Runtime results live only for this run. Public events remain the ordering signal.
type privateCompletions struct {
	sync.Mutex
	active, closed bool
	results        map[string]map[string]*model.Value
}

func withPrivateCompletions(ctx context.Context) context.Context {
	if privateResults(ctx) != nil {
		return ctx
	}
	return context.WithValue(ctx, privateCompletionKey{}, &privateCompletions{results: map[string]map[string]*model.Value{}})
}
func privateResults(ctx context.Context) *privateCompletions {
	results, _ := ctx.Value(privateCompletionKey{}).(*privateCompletions)
	return results
}
func activatePrivateResults(ctx context.Context) {
	if results := privateResults(ctx); results != nil {
		results.Lock()
		defer results.Unlock()
		results.active = true
	}
}
func recordPrivateCompletion(ctx context.Context, id string, outputs map[string]*model.Value) {
	if results := privateResults(ctx); results != nil {
		results.Lock()
		defer results.Unlock()
		if results.active && !results.closed {
			results.results[id] = outputs
		}
	}
}
func collectCompletion(ctx context.Context, end *model.CallEnded) map[string]*model.Value {
	if results := privateResults(ctx); results != nil {
		results.Lock()
		defer results.Unlock()
		if outputs, ok := results.results[end.Call.ID]; ok {
			delete(results.results, end.Call.ID)
			return outputs
		}
	}
	return end.Outputs
}
func releaseCompletion(ctx context.Context, id string) {
	if results := privateResults(ctx); results != nil {
		results.Lock()
		defer results.Unlock()
		delete(results.results, id)
	}
}
func clearPrivateCompletions(ctx context.Context) {
	if results := privateResults(ctx); results != nil {
		results.Lock()
		defer results.Unlock()
		results.closed = true
		clear(results.results)
	}
}

func privateResultsActive(ctx context.Context) bool {
	if results := privateResults(ctx); results != nil {
		results.Lock()
		defer results.Unlock()
		return results.active
	}
	return false
}

type resultSubtreeKey struct{}
type resultSubtree struct {
	parent *resultSubtree
	active atomic.Bool
}

func withResultSubtree(ctx context.Context) context.Context {
	parent, _ := ctx.Value(resultSubtreeKey{}).(*resultSubtree)
	return context.WithValue(ctx, resultSubtreeKey{}, &resultSubtree{parent: parent})
}

// Activation reaches ancestor ops, but never a sibling op.
func activateResultSubtree(ctx context.Context) {
	subtree, _ := ctx.Value(resultSubtreeKey{}).(*resultSubtree)
	for ; subtree != nil; subtree = subtree.parent {
		subtree.active.Store(true)
	}
}

func resultSubtreeActive(ctx context.Context) bool {
	subtree, _ := ctx.Value(resultSubtreeKey{}).(*resultSubtree)
	return subtree != nil && subtree.active.Load()
}
