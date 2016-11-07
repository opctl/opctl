package models

type ContainerStdOutWrittenToEvent struct {
  RootOpRunId string `json:"rootOpRunId"`
  OpRunId     string `json:"opRunId"`
  Data        []byte `json:"data"`
}
