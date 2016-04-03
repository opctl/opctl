package models

func NewAddSubOperationReq(
projectUrl *ProjectUrl,
isOperationSubOperation bool,
subOperationName string,
operationName string,
precedingSubOperationName string,
) *AddSubOperationReq {

  return &AddSubOperationReq{
    ProjectUrl:projectUrl,
    IsOperationSubOperation:isOperationSubOperation,
    SubOperationName :subOperationName,
    OperationName :operationName,
    PrecedingSubOperationName :precedingSubOperationName,
  }

}

type AddSubOperationReq struct {
  ProjectUrl                *ProjectUrl
  IsOperationSubOperation   bool
  SubOperationName          string
  OperationName             string
  PrecedingSubOperationName string
}