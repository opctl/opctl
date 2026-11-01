package node

import (
	"context"
	"time"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/goop"
)

type describeMemo struct {
	data []byte
	used uint64
}

func (cache *goBootstrapCache) cachedDescribe(key string) ([]byte, bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	entry, ok := cache.describes[key]
	if !ok {
		return nil, false
	}
	cache.describeClock++
	entry.used = cache.describeClock
	cache.describes[key] = entry
	return append([]byte(nil), entry.data...), true
}
func (cache *goBootstrapCache) rememberDescribe(key string, data []byte) {
	const budget = 64 * 1024 * 1024
	if len(data) > budget {
		return
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()
	if cache.describes == nil {
		cache.describes = map[string]describeMemo{}
	}
	cache.describeBytes -= len(cache.describes[key].data)
	cache.describeClock++
	cache.describes[key] = describeMemo{append([]byte(nil), data...), cache.describeClock}
	cache.describeBytes += len(data)
	for len(cache.describes) > 128 || cache.describeBytes > budget {
		var oldest string
		var used uint64
		for key, entry := range cache.describes {
			if used == 0 || entry.used < used {
				oldest, used = key, entry.used
			}
		}
		cache.describeBytes -= len(cache.describes[oldest].data)
		delete(cache.describes, oldest)
	}
}
func (bootstrap _bootstrapper) describeCached(ctx context.Context, op *model.OpCall, root string, program *goProgram, invoke func() ([]byte, error)) ([]byte, error) {
	unlock, err := bootstrap.cache.acquire(ctx, "describe:"+program.key)
	if err != nil {
		return nil, err
	}
	defer unlock()
	if data, ok := bootstrap.cache.cachedDescribe(program.key); ok {
		if err := bootstrap.cachedDescribeEvents(op, root); err != nil {
			return nil, err
		}
		return data, nil
	}
	data, err := invoke()
	if err != nil {
		return nil, err
	}
	if _, err := goop.ParseDescribe(data); err != nil {
		return nil, &OpDeclarationError{Err: err}
	}
	bootstrap.cache.rememberDescribe(program.key, data)
	return data, nil
}
func (bootstrap _bootstrapper) cachedDescribeEvents(op *model.OpCall, root string) error {
	if bootstrap.cache.pubSub == nil {
		return nil
	}
	id, err := uniquestring.Construct()
	if err != nil {
		return err
	}
	if root == "" {
		root = id
	}
	name := "bootstrap describe (cached)"
	call := model.Call{ID: id, RootID: root, Name: &name}
	if op.OpID != "" {
		call.ParentID = &op.OpID
	}
	bootstrap.cache.pubSub.Publish(model.Event{Timestamp: time.Now().UTC(), CallStarted: &model.CallStarted{Call: call, Ref: op.OpPath}})
	bootstrap.cache.pubSub.Publish(model.Event{Timestamp: time.Now().UTC(), CallEnded: &model.CallEnded{Call: call, Ref: op.OpPath, Outcome: model.OpOutcomeSucceeded, Outputs: map[string]*model.Value{}}})
	return nil
}
