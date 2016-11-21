package model

// Represents an op run. From a data structure perspective, op runs are directed rooted trees (aka out-tree's)
type OpRun struct {
  Id       string `json:"id"`
  OpRef    string `json:"opRef"`
  ParentId string `json:"parentId"`
  RootId   string `json:"rootId"`
}
