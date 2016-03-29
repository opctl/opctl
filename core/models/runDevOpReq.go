package models

func NewRunDevOpReq(
pathToProjectRootDir string,
devOpName string,
) *RunDevOpReq {

  return &RunDevOpReq{
    PathToProjectRootDir:pathToProjectRootDir,
    DevOpName :devOpName,
  }

}

type RunDevOpReq struct {
  PathToProjectRootDir string
  DevOpName            string
}
