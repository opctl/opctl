package models

type ContainerStdOutWrittenToEvent struct {
  OpRunDescriptor *OpRunDescriptor `json:"opRunDescriptor"`
  Data        []byte `json:"data"`
}
