package core

import (
	"sync"

	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o internal/fakes/callStore.go . callStore
// stores Call (dynamic call graph) calls
type callStore interface {
	// adds the provided call
	Add(call *model.Call)
	// lists all calls w/ parentID
	ListWithParentID(parentID string) []*model.Call
	// sets IsKilled to true
	SetIsKilled(id string)
	TryGet(id string) *model.Call
}

func newCallStore() callStore {

	return &_callStore{
		callsByID: make(map[string]*model.Call),
	}

}

type _callStore struct {
	callsByID map[string]*model.Call
	// synchronize access via mutex
	mux sync.RWMutex
}

// O(1) complexity
func (cs *_callStore) Add(call *model.Call) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	cs.callsByID[call.Id] = call
}

// O(n) complexity (n being active call count)
func (cs *_callStore) ListWithParentID(parentID string) []*model.Call {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	results := []*model.Call{}
	for _, call := range cs.callsByID {
		if nil != call.ParentID && *call.ParentID == parentID {
			results = append(results, call)
		}
	}
	return results
}

func (cs *_callStore) SetIsKilled(id string) {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	if _, ok := cs.callsByID[id]; ok {
		cs.callsByID[id].IsKilled = true
	}
}

func (cs *_callStore) TryGet(
	id string,
) *model.Call {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	if dcg, ok := cs.callsByID[id]; ok {
		return dcg
	}

	return nil
}
