package model

// Represents an op instance. From a data structure perspective, op instances are nodes in directed rooted tree (aka out-tree's)
type OpInstance struct {
  Id       string `json:"id"`
  OpRef    string `json:"opRef"`
  ParentId string `json:"parentId"`
  RootId   string `json:"rootId"`
}
