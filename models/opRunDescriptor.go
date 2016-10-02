package models

type OpRunDescriptor struct {
  Id       string `json:"id"`
  Op    string `json:"op"`
  Parent   *OpRunDescriptor `json:"parent"`
  Children []*OpRunDescriptor `json:"-"`
}
