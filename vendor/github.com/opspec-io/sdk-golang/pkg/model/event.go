package model

import "time"

type Event struct {
	ContainerExited          *ContainerExitedEvent          `json:"containerExitedEvent,omitempty"`
	ContainerStarted         *ContainerStartedEvent         `json:"containerStartedEvent,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitEmpty"`
	OpEnded                  *OpEndedEvent                  `json:"opEnded,omitempty"`
	OpStarted                *OpStartedEvent                `json:"opStarted,omitempty"`
	OpEncounteredError       *OpEncounteredErrorEvent       `json:"opEncounteredError,omitempty"`
	Timestamp                time.Time                      `json:"timestamp"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
)

type ContainerExitedEvent struct {
	ContainerRef string `json:"containerRef"`
	ExitCode     int    `json:"exitCode"`
	OpGraphId    string `json:"opGraphId"`
	ContainerId  string `json:"containerId"`
	OpRef        string `json:"opRef"`
}

type ContainerStartedEvent struct {
	ContainerRef string `json:"containerRef"`
	OpGraphId    string `json:"opGraphId"`
	ContainerId  string `json:"containerId"`
	OpRef        string `json:"opRef"`
}

type ContainerStdErrWrittenToEvent struct {
	ContainerRef string `json:"containerRef"`
	Data         []byte `json:"data"`
	OpGraphId    string `json:"opGraphId"`
	ContainerId  string `json:"containerId"`
	OpRef        string `json:"opRef"`
}

type ContainerStdOutWrittenToEvent struct {
	ContainerRef string `json:"containerRef"`
	Data         []byte `json:"data"`
	OpGraphId    string `json:"opGraphId"`
	ContainerId  string `json:"containerId"`
	OpRef        string `json:"opRef"`
}

type OpEncounteredErrorEvent struct {
	OpGraphId string `json:"opGraphId"`
	Msg       string `json:"msg"`
	OpId      string `json:"opId"`
	OpRef     string `json:"opRef"`
}

type OpEndedEvent struct {
	OpGraphId string `json:"opGraphId"`
	OpId      string `json:"opId"`
	OpRef     string `json:"opRef"`
	Outcome   string `json:"outcome"`
}

type OpStartedEvent struct {
	OpGraphId string `json:"opGraphId"`
	OpId      string `json:"opId"`
	OpRef     string `json:"opRef"`
}
