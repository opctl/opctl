package core

import (
	"sync"

	"github.com/opctl/opctl/sdks/go/types"
)

//go:generate counterfeiter -o ./fakeCallStore.go --fake-name fakeCallStore ./ callStore

// stores DCG (dynamic call graph) calls
type callStore interface {
	// adds the provided call
	Add(call *types.DCG)
	// lists all calls w/ parentID
	ListWithParentID(parentID string) []*types.DCG
	// sets IsKilled to true
	SetIsKilled(id string)
	Get(id string) types.DCG
}

func newCallStore() callStore {

	return &_callStore{
		callsByID: make(map[string]*types.DCG),
	}

}

type _callStore struct {
	callsByID map[string]*types.DCG
	// synchronize access via mutex
	mux sync.RWMutex
}

// O(1) complexity
func (cs *_callStore) Add(call *types.DCG) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	cs.callsByID[call.Id] = call
}

// O(n) complexity (n being active call count)
func (cs *_callStore) ListWithParentID(parentID string) []*types.DCG {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	results := []*types.DCG{}
	for _, call := range cs.callsByID {
		if nil != call.ParentID && *call.ParentID == parentID {
			results = append(results, call)
		}
	}
	return results
}

func (cs *_callStore) SetIsKilled(id string) {
	cs.callsByID[id].IsKilled = true
}

func (cs *_callStore) Get(
	id string,
) types.DCG {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	return *cs.callsByID[id]
}
