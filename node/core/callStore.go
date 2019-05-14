package core

import (
	"github.com/opctl/sdk-golang/model"
)

//go:generate counterfeiter -o ./fakeCallStore.go --fake-name fakeCallStore ./ callStore

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
	// synchronize access via mutex
	callsByID map[string]*model.DCG
}

// O(1) complexity
func (this *_callStore) Add(call *model.DCG) {
	this.callsByID[call.Id] = call
}

// O(n) complexity (n being active call count)
func (this *_callStore) ListWithParentID(parentID string) []*model.DCG {
	results := []*model.DCG{}

	for _, call := range this.callsByID {
		if nil != call.ParentID && *call.ParentID == parentID {
			results = append(results, call)
		}
	}
	return results
}

func (this *_callStore) SetIsKilled(id string) {
	this.callsByID[id].IsKilled = true
}

func (this _callStore) Get(
	id string,
) model.DCG {
	return *this.callsByID[id]
}
