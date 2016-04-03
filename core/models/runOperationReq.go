package models

func NewRunOperationReq(
projectUrl *ProjectUrl,
operationName string,
) *RunOperationReq {

  return &RunOperationReq{
    ProjectUrl:projectUrl,
    OperationName :operationName,
  }

}

type RunOperationReq struct {
  ProjectUrl    *ProjectUrl
  OperationName string
}
