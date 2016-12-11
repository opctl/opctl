package model

import "time"

type Event struct {
  ContainerStdErrWrittenTo      *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
  ContainerStdOutWrittenTo      *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitEmpty"`
  OpInstanceFailed              *OpInstanceEndedEvent `json:"opInstanceFailed,omitempty"`
  OpInstanceKilledEvent         *OpInstanceKilledEvent `json:"opInstanceKilled,omitempty"`
  OpInstanceStarted             *OpInstanceStartedEvent `json:"opInstanceStarted,omitempty"`
  OpInstanceEncounteredError    *OpInstanceEncounteredErrorEvent `json:"opInstanceEncounteredError,omitempty"`
  Timestamp                     time.Time `json:"timestamp"`
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

type OpInstanceEncounteredErrorEvent struct {
  Msg              string `json:"msg"`
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}

type OpInstanceEndedEvent struct {
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  Outcome          string `json:"outcome"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}

type OpInstanceKilledEvent struct {
  OpInstanceId string `json:"opInstanceId"`
}

type OpInstanceStartedEvent struct {
  OpRef            string `json:"opRef"`
  OpInstanceId     string `json:"opInstanceId"`
  RootOpInstanceId string `json:"rootOpInstanceId"`
}
