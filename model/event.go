package model

import "time"

// Event represents a distributed state change
type Event struct {
	ContainerExited          *ContainerExitedEvent          `json:"containerExited,omitempty"`
	ContainerStarted         *ContainerStartedEvent         `json:"containerStarted,omitempty"`
	ContainerStdErrWrittenTo *ContainerStdErrWrittenToEvent `json:"containerStdErrWrittenTo,omitempty"`
	ContainerStdOutWrittenTo *ContainerStdOutWrittenToEvent `json:"containerStdOutWrittenTo,omitempty"`
	OpEnded                  *OpEndedEvent                  `json:"opEnded,omitempty"`
	OpStarted                *OpStartedEvent                `json:"opStarted,omitempty"`
	OpEncounteredError       *OpEncounteredErrorEvent       `json:"opEncounteredError,omitempty"`
	Timestamp                time.Time                      `json:"timestamp"`
	OutputInitialized        *OutputInitializedEvent        `json:"outputInitialized,omitempty"`
	ParallelCallEnded        *ParallelCallEndedEvent        `json:"parallelCallEnded,omitempty"`
	SerialCallEnded          *SerialCallEndedEvent          `json:"serialCallEnded,omitempty"`
}

const (
	OpOutcomeSucceeded = "SUCCEEDED"
	OpOutcomeFailed    = "FAILED"
	OpOutcomeKilled    = "KILLED"
)

// ContainerExitedEvent represents the exit of a containerized process; no further events will occur for the container
type ContainerExitedEvent struct {
	ImageRef    string `json:"imageRef"`
	ExitCode    int    `json:"exitCode"`
	RootOpId    string `json:"rootOpId"`
	ContainerId string `json:"containerId"`
	PkgRef      string `json:"pkgRef"`
}

type ContainerStartedEvent struct {
	ImageRef    string `json:"imageRef"`
	RootOpId    string `json:"rootOpId"`
	ContainerId string `json:"containerId"`
	PkgRef      string `json:"pkgRef"`
}

// ContainerStdErrWrittenToEvent represents a single write to a containers std err.
type ContainerStdErrWrittenToEvent struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootOpId    string `json:"rootOpId"`
	ContainerId string `json:"containerId"`
	PkgRef      string `json:"pkgRef"`
}

// ContainerStdOutWrittenToEvent represents a single write to a containers std out.
type ContainerStdOutWrittenToEvent struct {
	ImageRef    string `json:"imageRef"`
	Data        []byte `json:"data"`
	RootOpId    string `json:"rootOpId"`
	ContainerId string `json:"containerId"`
	PkgRef      string `json:"pkgRef"`
}

type OpEncounteredErrorEvent struct {
	RootOpId string `json:"rootOpId"`
	Msg      string `json:"msg"`
	OpId     string `json:"opId"`
	PkgRef   string `json:"pkgRef"`
}

// OpEndedEvent represents the end of an op; no further events will occur for the op.
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

type OutputInitializedEvent struct {
	Name     string `json:"name"`
	CallId   string `json:"callId"`
	Value    *Data  `json:"value"`
	RootOpId string `json:"rootOpId"`
}

type ParallelCallEndedEvent struct {
	RootOpId string `json:"rootOpId"`
	CallId   string `json:"callId"`
}

type SerialCallEndedEvent struct {
	RootOpId string `json:"rootOpId"`
	CallId   string `json:"callId"`
}
