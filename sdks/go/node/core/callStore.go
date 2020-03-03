package core

import (
	"sync"

	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o internal/fakes/callStore.go . callStore
// stores DCG (dynamic call graph) calls
type callStore interface {
	// adds the provided call
	Add(call *model.DCG)
	// lists all calls w/ parentID
	ListWithParentID(parentID string) []*model.DCG
	// sets IsKilled to true
	SetIsKilled(id string)
	Get(id string) model.DCG
}

func newCallStore() callStore {

	return &_callStore{
		callsByID: make(map[string]*model.DCG),
	}

}

type _callStore struct {
	callsByID map[string]*model.DCG
	// synchronize access via mutex
	mux sync.RWMutex
}

// O(1) complexity
func (cs *_callStore) Add(call *model.DCG) {
	cs.mux.Lock()
	defer cs.mux.Unlock()

	cs.callsByID[call.Id] = call
}

// O(n) complexity (n being active call count)
func (cs *_callStore) ListWithParentID(parentID string) []*model.DCG {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	results := []*model.DCG{}
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
) model.DCG {
	cs.mux.RLock()
	defer cs.mux.RUnlock()

	return *cs.callsByID[id]
}
