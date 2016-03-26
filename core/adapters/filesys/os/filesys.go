package os

import "github.com/dev-op-spec/engine/core/ports"

func New() ports.Filesys {

  return filesys{
    compositionRoot:newCompositionRoot(),
  }

}

type filesys struct {
  compositionRoot compositionRoot
}

func (this filesys) CreateDevOpDir(
devOpName string,
) (err error) {
  return this.compositionRoot.
  CreateDevOpDirUseCase().
  Execute(devOpName)
}

func (this filesys) CreatePipelineDir(
pipelineName string,
) (err error) {
  return this.compositionRoot.
  CreatePipelineDirUseCase().
  Execute(pipelineName)
}

func (this filesys) ListNamesOfDevOpDirs(
) (namesOfDevOpDirs []string, err error) {
  return this.compositionRoot.
  ListNamesOfDevOpDirsUseCase().
  Execute()
}

func (this filesys) ListNamesOfPipelineDirs(
) (namesOfPipelineDirs []string, err error) {
  return this.compositionRoot.
  ListNamesOfPipelineDirsUseCase().
  Execute()
}

func (this filesys) ReadDevOpFile(
devOpName string,
) (file []byte, err error) {
  return this.compositionRoot.
  ReadDevOpFileUseCase().
  Execute(devOpName)
}

func (this filesys) ReadPipelineFile(
pipelineName string,
) (file []byte, err error) {
  return this.compositionRoot.
  ReadPipelineFileUseCase().
  Execute(pipelineName)
}

func (this filesys) SaveDevOpFile(
devOpName string,
data []byte,
) (err error) {
  return this.compositionRoot.
  SaveDevOpFileUseCase().
  Execute(
    devOpName,
    data,
  )
}

func (this filesys) SavePipelineFile(
pipelineName string,
data []byte,
) (err error) {
  return this.compositionRoot.
  SavePipelineFileUseCase().
  Execute(
    pipelineName,
    data,
  )
}

const (
  relPathToDevOpSpecDir = "./.dev-op-spec"

  relPathToDevOpsDir = relPathToDevOpSpecDir + "/dev-ops"

  relPathToPipelinesDir = relPathToDevOpSpecDir + "/pipelines"
)
