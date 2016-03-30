package models

func NewRunDevOpReq(
projectUrl *ProjectUrl,
devOpName string,
) *RunDevOpReq {

  return &RunDevOpReq{
    ProjectUrl:projectUrl,
    DevOpName :devOpName,
  }

}

type RunDevOpReq struct {
  ProjectUrl *ProjectUrl
  DevOpName  string
}
