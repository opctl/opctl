package osfilesys

import "github.com/dev-op-spec/engine/core/ports"

func NewFilesys() ports.Filesys {

  return filesys{
    compositionRoot:newCompositionRoot(),
  }

}

type filesys struct {
  compositionRoot compositionRoot
}

func (fs filesys) CreateDevOpDir(
devOpName string,
) (err error) {
  return fs.compositionRoot.
  CreateDevOpDirUcExecuter().
  Execute(devOpName)
}

func (fs filesys) CreatePipelineDir(
pipelineName string,
) (err error) {
  return fs.compositionRoot.
  CreatePipelineDirUcExecuter().
  Execute(pipelineName)
}

func (fs filesys) ListNamesOfDevOpDirs(
) (namesOfDevOpDirs []string, err error) {
  return fs.compositionRoot.
  ListNamesOfDevOpDirsUcExecuter().
  Execute()
}

func (fs filesys) ListNamesOfPipelineDirs(
) (namesOfPipelineDirs []string, err error) {
  return fs.compositionRoot.
  ListNamesOfPipelineDirsUcExecuter().
  Execute()
}

func (fs filesys) ReadDevOpFile(
devOpName string,
) (file []byte, err error) {
  return fs.compositionRoot.
  ReadDevOpFileUcExecuter().
  Execute(devOpName)
}

func (fs filesys) ReadPipelineFile(
pipelineName string,
) (file []byte, err error) {
  return fs.compositionRoot.
  ReadPipelineFileUcExecuter().
  Execute(pipelineName)
}

func (fs filesys) SaveDevOpFile(
devOpName string,
data []byte,
) (err error) {
  return fs.compositionRoot.
  SaveDevOpFileUcExecuter().
  Execute(
    devOpName,
    data,
  )
}

func (fs filesys) SavePipelineFile(
pipelineName string,
data []byte,
) (err error) {
  return fs.compositionRoot.
  SavePipelineFileUcExecuter().
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
