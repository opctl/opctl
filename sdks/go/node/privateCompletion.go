package node

import (
	"context"
	"sync"

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
