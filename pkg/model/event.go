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
	OutputInitialized        *OutputInitializedEvent        `json:"outputInitialized,omitempty"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
)

type OutputInitializedEvent struct {
	Name     string `json:"name"`
	CallId   string `json:"callId"`
	Value    *Data  `json:"value"`
	RootOpId string `json:"rootOpId"`
}

type ContainerExitedEvent struct {
	ContainerRef string `json:"containerRef"`
	ExitCode     int    `json:"exitCode"`
	RootOpId     string `json:"rootOpId"`
	ContainerId  string `json:"containerId"`
	PkgRef       string `json:"pkgRef"`
}

type ContainerStartedEvent struct {
	ContainerRef string `json:"containerRef"`
	RootOpId     string `json:"rootOpId"`
	ContainerId  string `json:"containerId"`
	PkgRef       string `json:"pkgRef"`
}

type ContainerStdErrWrittenToEvent struct {
	ContainerRef string `json:"containerRef"`
	Data         []byte `json:"data"`
	RootOpId     string `json:"rootOpId"`
	ContainerId  string `json:"containerId"`
	PkgRef       string `json:"pkgRef"`
}

type ContainerStdOutWrittenToEvent struct {
	ContainerRef string `json:"containerRef"`
	Data         []byte `json:"data"`
	RootOpId     string `json:"rootOpId"`
	ContainerId  string `json:"containerId"`
	PkgRef       string `json:"pkgRef"`
}

type OpEncounteredErrorEvent struct {
	RootOpId string `json:"rootOpId"`
	Msg      string `json:"msg"`
	OpId     string `json:"opId"`
	PkgRef   string `json:"pkgRef"`
}

type OpEndedEvent struct {
	RootOpId string `json:"rootOpId"`
	OpId     string `json:"opId"`
	PkgRef   string `json:"pkgRef"`
	Outcome  string `json:"outcome"`
}

type OpStartedEvent struct {
	RootOpId string `json:"rootOpId"`
	OpId     string `json:"opId"`
	PkgRef   string `json:"pkgRef"`
}
