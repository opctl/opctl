package core

//go:generate counterfeiter -o ./fakeNodeRepo.go --fake-name fakeNodeRepo ./ nodeRepo

import (
	"sync"
)

type nodeRepo interface {
	// adds the provided node
	add(node *nodeDescriptor)
	// deletes the node with the provided id
	deleteIfExists(nodeId string)
	// lists all nodes with the provided opGraphId
	listWithOpGraphId(opGraphId string) []*nodeDescriptor
	// tries to get the node with the provided id; returns nil if not found
	getIfExists(nodeId string) *nodeDescriptor
}

func newNodeRepo() nodeRepo {

	return &_nodeRepo{
		byIdIndex:      make(map[string]*nodeDescriptor),
		byIdIndexMutex: sync.RWMutex{},
	}

}

type _nodeRepo struct {
	// synchronize access via mutex
	byIdIndex      map[string]*nodeDescriptor
	byIdIndexMutex sync.RWMutex
}

// O(1) complexity; thread safe
func (this *_nodeRepo) add(node *nodeDescriptor) {
	this.byIdIndexMutex.Lock()
	defer this.byIdIndexMutex.Unlock()
	this.byIdIndex[node.Id] = node
}

// O(1) complexity; thread safe
func (this *_nodeRepo) deleteIfExists(nodeId string) {
	this.byIdIndexMutex.Lock()
	defer this.byIdIndexMutex.Unlock()

	delete(this.byIdIndex, nodeId)
}

// O(n) complexity (n being active node count); thread safe
func (this *_nodeRepo) listWithOpGraphId(opGraphId string) []*nodeDescriptor {
	this.byIdIndexMutex.RLock()
	defer this.byIdIndexMutex.RUnlock()

	nodesWithGraphIdSlice := []*nodeDescriptor{}

	for _, node := range this.byIdIndex {
		if node.OpGraphId == opGraphId {
			nodesWithGraphIdSlice = append(nodesWithGraphIdSlice, node)
		}
	}
	return nodesWithGraphIdSlice
}

// O(1) complexity; thread safe
func (this *_nodeRepo) getIfExists(nodeId string) *nodeDescriptor {
	this.byIdIndexMutex.RLock()
	defer this.byIdIndexMutex.RUnlock()

	return this.byIdIndex[nodeId]
}
