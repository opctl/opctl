package models

type OpRunStartedEvent struct {
  Descriptor *OpRunDescriptor `json:"descriptor"`
}
