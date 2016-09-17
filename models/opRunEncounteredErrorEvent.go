package models

type OpRunEncounteredErrorEvent struct {
  OpRunId       string `json:"opRunId"`
  OpRef         string `json:"opRef"`
  Msg           string `json:"msg"`
  ParentOpRunId string `json:"parentOpRunId"`
  RootOpRunId   string `json:"rootOpRunId"`
}
