package models

func NewRunOperationReq(
operationUrl *Url,
) *RunOperationReq {

  return &RunOperationReq{
    OperationUrl:operationUrl,
  }

}

type RunOperationReq struct {
  OperationUrl  *Url
}
