package models

func NewAddSubOpReq(
projectUrl *Url,
opName string,
subOpUrl string,
precedingSubOpUrl string,
) *AddSubOpReq {

  return &AddSubOpReq{
    ProjectUrl:projectUrl,
    OpName :opName,
    SubOpUrl :subOpUrl,
    PrecedingSubOpUrl :precedingSubOpUrl,
  }

}

type AddSubOpReq struct {
  ProjectUrl               *Url
  OpName            string `json:"opName"`
  SubOpUrl          string `json:"subOpUrl"`
  PrecedingSubOpUrl string `json:"precedingSubOpUrl"`
}