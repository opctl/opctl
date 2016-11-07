package models

// Represents an op run. From a data structure perspective, op runs are directed rooted trees (aka out-tree's)
type OpRun struct {
  Children []string `json:"children,omitempty"`
  Id       string `json:"id"`
  Parent   string `json:"parent,omitempty"`
  OpRef    string `json:"opRef"`
}
