package model

// Represents an op run. From a data structure perspective, op runs are directed rooted trees (aka out-tree's)
type OpRun struct {
  Children []*OpRun `json:"children,omitempty"`
  Id       string `json:"id"`
  ParentId string `json:"-"`
  OpRef    string `json:"opRef"`
}
