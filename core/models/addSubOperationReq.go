package models

func NewAddSubOperationReq(
projectUrl *ProjectUrl,
subOperationName string,
operationName string,
precedingSubOperationName string,
) *AddSubOperationReq {

  return &AddSubOperationReq{
    ProjectUrl:projectUrl,
    SubOperationName :subOperationName,
    OperationName :operationName,
    PrecedingSubOperationName :precedingSubOperationName,
  }

}

type AddSubOperationReq struct {
  ProjectUrl                *ProjectUrl
  SubOperationName          string `json:"subOperationName"`
  OperationName             string `json:"operationName"`
  PrecedingSubOperationName string `json:"precedingSubOperationName"`
}