package models

type OpRunEncounteredErrorEvent struct {
  Descriptor *OpRunDescriptor `json:"descriptor"`
  Msg        string `json:"msg"`
}
