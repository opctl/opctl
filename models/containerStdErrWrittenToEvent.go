package models

type ContainerStdErrWrittenToEvent struct {
  OpRunDescriptor *OpRunDescriptor `json:"opRunDescriptor"`
  Data        []byte `json:"data"`
}
