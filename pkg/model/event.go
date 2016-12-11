package model

import "time"

type Event struct {
  ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
  ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitEmpty"`
  OpEnded                  *OpEndedEvent `json:"opEnded,omitempty"`
  OpStarted                *OpStartedEvent `json:"opStarted,omitempty"`
  OpEncounteredError       *OpEncounteredErrorEvent `json:"opEncounteredError,omitempty"`
  Timestamp                time.Time `json:"timestamp"`
}

const (
  OpOutcomeSucceeded = "SUCCEEDED"
  OpOutcomeFailed = "FAILED"
  OpOutcomeKilled = "KILLED"
)

type ContainerStdErrWrittenToEvent struct {
  Data             []byte `json:"data"`
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}

type ContainerStdOutWrittenToEvent struct {
  Data             []byte `json:"data"`
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}

type OpEncounteredErrorEvent struct {
  Msg              string `json:"msg"`
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}

type OpEndedEvent struct {
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  Outcome          string `json:"outcome"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}

type OpStartedEvent struct {
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}
