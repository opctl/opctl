package models

import "time"

type Event struct {
  Timestamp                time.Time `json:"timestamp"`
  OpRunEnded               *OpRunEndedEvent `json:"opRunEnded,omitempty"`
  OpRunStarted             *OpRunStartedEvent `json:"opRunStarted,omitempty"`
  OpRunEncounteredError    *OpRunEncounteredErrorEvent `json:"opRunEncounteredError,omitempty"`
  ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
  ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitEmpty"`
}
