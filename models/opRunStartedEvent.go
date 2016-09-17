package models

type OpRunStartedEvent struct {
  Id       string `json:"id"`
  OpRef    string `json:"opRef"`
  ParentId string `json:"parentId"`
  RootId   string `json:"rootId"`
}
