package model

// Represents an op instance. From a data structure perspective, op instances are nodes in directed rooted tree (aka out-tree's)
type OpInstanceView struct {
  Id       string `json:"id"`
  OpRef    string `json:"opRef"`
  Run []*CallGraphInstance `json:"run"`
  ParentId string `json:"parentId"`
  RootId   string `json:"rootId"`
}
