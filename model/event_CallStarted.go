package model

import "strings"

type CallStartedEvent struct {
	*CallEventBase
	Container *ContainerCallStartedEvent `json:"container,omitempty" yaml:"container,omitempty"`
	Op        *OpCallStartedEvent        `json:"op,omitempty" yaml:"op,omitempty"`
	Parallel  *ParallelCallStartedEvent  `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	Serial    *SerialCallStartedEvent    `json:"serial,omitempty" yaml:"serial,omitempty"`
}

// implement fmt.Stringer interface
func (cse CallStartedEvent) String() string {
	var pre string
	switch {
	case nil != cse.Container:
		pre = cse.Container.String()
	case nil != cse.Op:
		pre = cse.Op.String()
	case nil != cse.Parallel:
		pre = cse.Parallel.String()
	case nil != cse.Serial:
		pre = cse.Serial.String()
	}

	return strings.Join(
		[]string{
			pre,
			cse.CallEventBase.String(),
		},
		" ",
	)
}

type ContainerCallStartedEvent struct{}

// implement fmt.Stringer interface
func (ccse ContainerCallStartedEvent) String() string {
	return "ContainerCallStarted"
}

type OpCallStartedEvent struct{}

// implement fmt.Stringer interface
func (ocse OpCallStartedEvent) String() string {
	return "OpCallStarted"
}

type ParallelCallStartedEvent struct{}

// implement fmt.Stringer interface
func (pcse ParallelCallStartedEvent) String() string {
	return "ParallelCallStarted"
}

type SerialCallStartedEvent struct{}

// implement fmt.Stringer interface
func (scse SerialCallStartedEvent) String() string {
	return "SerialCallStarted"
}
