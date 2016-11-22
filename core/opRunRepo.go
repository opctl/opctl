package core

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
  "sync"
)

type opRunRepo interface {
  // adds the provided op
  add(opRun *model.OpRun)
  // deletes the op with the provided id
  delete(opRunId string)
  // lists all op runs with the provided parentId
  listWithParentId(parentId string) []*model.OpRun
  // tries to get the op with the provided id; returns nil if not found
  tryGet(opRunId string) *model.OpRun
}

func newOpRunRepo() opRunRepo {

  return &_opRunRepo{
    byIdIndex:make(map[string]*model.OpRun),
    byIdIndexMutex:sync.RWMutex{},
  }

}

type _opRunRepo struct {
  // maintain dual indexes and synchronize access via mutex
  byIdIndex      map[string]*model.OpRun
  byIdIndexMutex sync.RWMutex
}

// O(1) complexity; thread safe
func (this *_opRunRepo) add(opRun *model.OpRun) {
  this.byIdIndexMutex.Lock()
  defer this.byIdIndexMutex.Unlock()
  this.byIdIndex[opRun.Id] = opRun
}

// O(1) complexity; thread safe
func (this *_opRunRepo) delete(opRunId string) {
  this.byIdIndexMutex.Lock()
  defer this.byIdIndexMutex.Unlock()

  delete(this.byIdIndex, opRunId)
}

// O(n) complexity (n being running op run count); thread safe
func (this *_opRunRepo) listWithParentId(parentId string) []*model.OpRun {
  this.byIdIndexMutex.RLock()
  defer this.byIdIndexMutex.RUnlock()

  opRunsWithParentIdSlice := []*model.OpRun{}

  for _, opRun := range this.byIdIndex {
    if (opRun.ParentId == parentId) {
      opRunsWithParentIdSlice = append(opRunsWithParentIdSlice, opRun)
    }
  }
  return opRunsWithParentIdSlice
}

// O(1) complexity; thread safe
func (this *_opRunRepo) tryGet(opRunId string) *model.OpRun {
  this.byIdIndexMutex.RLock()
  defer this.byIdIndexMutex.RUnlock()

  return this.byIdIndex[opRunId]
}
