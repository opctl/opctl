package models

func NewRunOpReq(
opUrl *Url,
) *RunOpReq {

  return &RunOpReq{
    OpUrl:opUrl,
  }

}

type RunOpReq struct {
  OpUrl  *Url
}
