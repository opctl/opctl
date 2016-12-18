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

type KillOpReq struct {
  OpInstanceId string
}

type SetCollectionDescriptionReq struct {
  PathToCollection string
  Description      string
}

type SetOpDescriptionReq struct {
  PathToOp    string
  Description string
}

type StartOpReq struct {
  // map of args keyed by param name
  Args  map[string]interface{} `json:"args"`
  OpRef string `json:"opRef"`
}
