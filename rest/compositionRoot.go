package rest

import (
  "github.com/dev-op-spec/engine/core"
  "net/http"
)

type compositionRoot interface {
  AddOperationHandler() http.Handler
  AddSubOperationHandler() http.Handler
  ListOperationsHandler() http.Handler
  RunOperationHandler() http.Handler
  SetDescriptionOfOperationHandler() http.Handler
}

func newCompositionRoot(
coreApi core.Api,
) (compositionRoot compositionRoot) {

  compositionRoot = &_compositionRoot{
    addOperationHandler:newAddOperationHandler(coreApi),
    addSubOperationHandler:newAddSubOperationHandler(coreApi),
    listOperationsHandler:newListOperationsHandler(coreApi),
    runOperationHandler:newRunOperationHandler(coreApi),
    setDescriptionOfOperationHandler:newSetDescriptionOfOperationHandler(coreApi),
  }

  return

}

type _compositionRoot struct {
  addOperationHandler              http.Handler
  addSubOperationHandler           http.Handler
  listOperationsHandler            http.Handler
  runOperationHandler              http.Handler
  setDescriptionOfOperationHandler http.Handler
}

func (this _compositionRoot) AddOperationHandler(
) http.Handler {
  return this.addOperationHandler
}

func (this _compositionRoot) AddSubOperationHandler(
) http.Handler {
  return this.addSubOperationHandler
}

func (this _compositionRoot) ListOperationsHandler(
) http.Handler {
  return this.listOperationsHandler
}

func (this _compositionRoot) RunOperationHandler(
) http.Handler {
  return this.runOperationHandler
}

func (this _compositionRoot) SetDescriptionOfOperationHandler(
) http.Handler {
  return this.setDescriptionOfOperationHandler
}