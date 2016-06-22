package models

func NewRunOpReq(
opUrl string,
args map[string]string,
) *RunOpReq {

  return &RunOpReq{
    OpUrl:opUrl,
    Args:args,
  }

}

type RunOpReq struct {
  OpUrl string
  Args  map[string]string
}
