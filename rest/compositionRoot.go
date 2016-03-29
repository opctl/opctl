package rest

import (
  "github.com/dev-op-spec/engine/core"
  "net/http"
)

type compositionRoot interface {
  AddDevOpHandler() http.Handler
  AddPipelineHandler() http.Handler
  AddStageToPipelineHandler() http.Handler
  ListDevOpsHandler() http.Handler
  ListPipelinesHandler() http.Handler
  RunDevOpHandler() http.Handler
  RunPipelineHandler() http.Handler
  SetDescriptionOfDevOpHandler() http.Handler
  SetDescriptionOfPipelineHandler() http.Handler
}

func newCompositionRoot(
coreApi core.Api,
) (compositionRoot compositionRoot) {

  compositionRoot = &_compositionRoot{
    addDevOpHandler: newAddDevOpHandler(coreApi),
    addPipelineHandler:newAddPipelineHandler(coreApi),
    addStageToPipelineHandler:newAddStageToPipelineHandler(coreApi),
    listDevOpsHandler:newListDevOpsHandler(coreApi),
    listPipelinesHandler:newListPipelinesHandler(coreApi),
    runDevOpHandler:newRunDevOpHandler(coreApi),
    runPipelineHandler:newRunPipelineHandler(coreApi),
    setDescriptionOfDevOpHandler:newSetDescriptionOfDevOpHandler(coreApi),
    setDescriptionOfPipelineHandler:newSetDescriptionOfPipelineHandler(coreApi),
  }

  return

}

type _compositionRoot struct {
  addDevOpHandler                 http.Handler
  addPipelineHandler              http.Handler
  addStageToPipelineHandler       http.Handler
  listDevOpsHandler               http.Handler
  listPipelinesHandler            http.Handler
  runDevOpHandler                 http.Handler
  runPipelineHandler              http.Handler
  setDescriptionOfDevOpHandler    http.Handler
  setDescriptionOfPipelineHandler http.Handler
}

func (this _compositionRoot) AddDevOpHandler(
) http.Handler {
  return this.addDevOpHandler
}

func (this _compositionRoot) AddPipelineHandler(
) http.Handler {
  return this.addPipelineHandler
}

func (this _compositionRoot) AddStageToPipelineHandler(
) http.Handler {
  return this.addStageToPipelineHandler
}

func (this _compositionRoot) ListDevOpsHandler(
) http.Handler {
  return this.listDevOpsHandler
}

func (this _compositionRoot) ListPipelinesHandler(
) http.Handler {
  return this.listPipelinesHandler
}

func (this _compositionRoot) RunDevOpHandler(
) http.Handler {
  return this.runDevOpHandler
}

func (this _compositionRoot) RunPipelineHandler(
) http.Handler {
  return this.runPipelineHandler
}

func (this _compositionRoot) SetDescriptionOfDevOpHandler(
) http.Handler {
  return this.setDescriptionOfDevOpHandler
}

func (this _compositionRoot) SetDescriptionOfPipelineHandler(
) http.Handler {
  return this.setDescriptionOfPipelineHandler
}