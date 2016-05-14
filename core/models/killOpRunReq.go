package models

func NewKillOpRunReq(
opRunId string,
) *KillOpRunReq {

  return &KillOpRunReq{
    OpRunId:opRunId,
  }

}

type KillOpRunReq struct {
  OpRunId string
}
