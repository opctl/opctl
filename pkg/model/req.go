package model

type CreateCollectionReq struct {
  Path        string
  Name        string
  Description string
}

type CreateOpReq struct {
  Path        string
  Name        string
  Description string
}

type KillOpRunReq struct {
  OpRunId string
}

type SetCollectionDescriptionReq struct {
PathToCollection string
Description  string
}

type SetOpDescriptionReq struct {
  PathToOp string
  Description  string
}

type StartOpRunReq struct {
  // map of args keyed by param name
  Args  map[string]*Arg `json:"args"`
  OpUrl string `json:"opUrl"`
}
