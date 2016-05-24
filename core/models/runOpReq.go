package models

func NewRunOpReq(
opUrl *Url,
args map[string]string,
) *RunOpReq {

  return &RunOpReq{
    OpUrl:opUrl,
    Args:args,
  }

}

type RunOpReq struct {
  OpUrl *Url
  Args  map[string]string
}
