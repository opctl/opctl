package models

func NewAddSubOperationReq(
projectUrl *Url,
operationName string,
subOperationUrl string,
precedingSubOperationUrl string,
) *AddSubOperationReq {

  return &AddSubOperationReq{
    ProjectUrl:projectUrl,
    OperationName :operationName,
    SubOperationUrl :subOperationUrl,
    PrecedingSubOperationUrl :precedingSubOperationUrl,
  }

}

type AddSubOperationReq struct {
  ProjectUrl               *Url
  OperationName            string `json:"operationName"`
  SubOperationUrl          string `json:"subOperationUrl"`
  PrecedingSubOperationUrl string `json:"precedingSubOperationUrl"`
}