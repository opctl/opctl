package model

import "time"

type Event struct {
  Timestamp                time.Time `json:"timestamp"`
  OpRunEnded               *OpRunEndedEvent `json:"opRunEnded,omitempty"`
  OpRunStarted             *OpRunStartedEvent `json:"opRunStarted,omitempty"`
  OpRunEncounteredError    *OpRunEncounteredErrorEvent `json:"opRunEncounteredError,omitempty"`
  ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
  ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitEmpty"`
}

const (
  OpRunOutcomeSucceeded = "SUCCEEDED"
  OpRunOutcomeFailed = "FAILED"
  OpRunOutcomeKilled = "KILLED"
)

type OpRunEndedEvent struct {
  OpRef         string `json:"opRef"`
  OpRunId       string `json:"opRunId"`
  Outcome       string `json:"outcome"`
  RootOpRunId   string `json:"rootOpRunId"`
}

type OpRunStartedEvent struct {
  OpRef         string `json:"opRef"`
  OpRunId       string `json:"opRunId"`
  RootOpRunId   string `json:"rootOpRunId"`
}

type OpRunEncounteredErrorEvent struct {
  Msg           string `json:"msg"`
  OpRef         string `json:"opRef"`
  OpRunId       string `json:"opRunId"`
  RootOpRunId   string `json:"rootOpRunId"`
}

type ContainerStdErrWrittenToEvent struct {
  Data        []byte `json:"data"`
  OpRef         string `json:"opRef"`
  OpRunId     string `json:"opRunId"`
  RootOpRunId string `json:"rootOpRunId"`
}

type ContainerStdOutWrittenToEvent struct {
  Data        []byte `json:"data"`
  OpRef         string `json:"opRef"`
  OpRunId     string `json:"opRunId"`
  RootOpRunId string `json:"rootOpRunId"`
}
