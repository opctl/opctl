package models

type ContainerStdErrWrittenToEvent struct {
  RootOpRunId string `json:"rootOpRunId"`
  OpRunId     string `json:"opRunId"`
  Data        []byte `json:"data"`
}
