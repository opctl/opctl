package core

import (
"github.com/dev-op-spec/engine/core/adapters/osfilesys"
"github.com/dev-op-spec/engine/core/adapters/dockercompose"
)

type compositionRoot interface {
  AddDevOpUcExecuter() addDevOpUcExecuter
  AddPipelineUcExecuter() addPipelineUcExecuter
  AddStageToPipelineUcExecuter() addStageToPipelineUcExecuter
  ListDevOpsUcExecuter() listDevOpsUcExecuter
  ListPipelinesUcExecuter() listPipelinesUcExecuter
  RunDevOpUcExecuter() runDevOpUcExecuter
  RunPipelineUcExecuter() runPipelineUcExecuter
  SetDescriptionOfDevOpUcExecuter() setDescriptionOfDevOpUcExecuter
  SetDescriptionOfPipelineUcExecuter() setDescriptionOfPipelineUcExecuter
}

func newCompositionRoot(
) (compositionRoot compositionRoot, err error) {

  fs := osfilesys.NewFilesys()

  yml := yamlCodecImpl{}

  containerEngine, err := dockercompose.NewContainerEngine()
  if (nil != err) {
    return
  }

  runDevOpUcExecuter := newRunDevOpUcExecuter(containerEngine)

  compositionRoot = &_compositionRoot{
    addDevOpUcExecuter: newAddDevOpUcExecuter(fs, yml, containerEngine),
    addPipelineUcExecuter: newAddPipelineUcExecuter(fs, yml),
    addStageToPipelineUcExecuter: newAddStageToPipelineUcExecuter(fs, yml),
    listDevOpsUcExecuter: newListDevOpsUcExecuter(fs, yml),
    listPipelinesUcExecuter: newListPipelinesUcExecuter(fs, yml),
    runDevOpUcExecuter: runDevOpUcExecuter,
    runPipelineUcExecuter: newRunPipelineUcExecuter(fs, yml, runDevOpUcExecuter),
    setDescriptionOfDevOpUcExecuter: newSetDescriptionOfDevOpUcExecuter(fs, yml),
    setDescriptionOfPipelineUcExecuter: newSetDescriptionOfPipelineUcExecuter(fs, yml),
  }

  return

}

type _compositionRoot struct {
  addDevOpUcExecuter                 addDevOpUcExecuter
  addPipelineUcExecuter              addPipelineUcExecuter
  addStageToPipelineUcExecuter       addStageToPipelineUcExecuter
  listDevOpsUcExecuter               listDevOpsUcExecuter
  listPipelinesUcExecuter            listPipelinesUcExecuter
  runDevOpUcExecuter                 runDevOpUcExecuter
  runPipelineUcExecuter              runPipelineUcExecuter
  setDescriptionOfDevOpUcExecuter    setDescriptionOfDevOpUcExecuter
  setDescriptionOfPipelineUcExecuter setDescriptionOfPipelineUcExecuter
}

func (_compositionRoot _compositionRoot) AddDevOpUcExecuter() addDevOpUcExecuter {
  return _compositionRoot.addDevOpUcExecuter
}

func (_compositionRoot _compositionRoot) AddPipelineUcExecuter() addPipelineUcExecuter {
  return _compositionRoot.addPipelineUcExecuter
}

func (_compositionRoot _compositionRoot) AddStageToPipelineUcExecuter() addStageToPipelineUcExecuter {
  return _compositionRoot.addStageToPipelineUcExecuter
}

func (_compositionRoot _compositionRoot) ListDevOpsUcExecuter() listDevOpsUcExecuter {
  return _compositionRoot.listDevOpsUcExecuter
}

func (_compositionRoot _compositionRoot) ListPipelinesUcExecuter() listPipelinesUcExecuter {
  return _compositionRoot.listPipelinesUcExecuter
}

func (_compositionRoot _compositionRoot) RunDevOpUcExecuter() runDevOpUcExecuter {
  return _compositionRoot.runDevOpUcExecuter
}

func (_compositionRoot _compositionRoot) RunPipelineUcExecuter() runPipelineUcExecuter {
  return _compositionRoot.runPipelineUcExecuter
}

func (_compositionRoot _compositionRoot) SetDescriptionOfDevOpUcExecuter() setDescriptionOfDevOpUcExecuter {
  return _compositionRoot.setDescriptionOfDevOpUcExecuter
}

func (_compositionRoot _compositionRoot) SetDescriptionOfPipelineUcExecuter() setDescriptionOfPipelineUcExecuter {
  return _compositionRoot.setDescriptionOfPipelineUcExecuter
}
