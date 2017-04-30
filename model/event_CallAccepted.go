package model

import "strings"

// CallAcceptedEvent represents acceptance of a call request.
type CallAcceptedEvent struct {
	*CallEventBase
	Container *ContainerCallAcceptedEvent `json:"container,omitempty" yaml:"container,omitempty"`
	Op        *OpCallAcceptedEvent        `json:"op,omitempty" yaml:"op,omitempty"`
	Parallel  *ParallelCallAcceptedEvent  `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	Serial    *SerialCallAcceptedEvent    `json:"serial,omitempty" yaml:"serial,omitempty"`
}

// implement fmt.Stringer interface
func (cae CallAcceptedEvent) String() string {
	var pre string
	switch {
	case nil != cae.Container:
		pre = cae.Container.String()
	case nil != cae.Op:
		pre = cae.Op.String()
	case nil != cae.Parallel:
		pre = cae.Parallel.String()
	case nil != cae.Serial:
		pre = cae.Serial.String()
	}

	return strings.Join(
		[]string{
			pre,
			cae.CallEventBase.String(),
		},
		" ",
	)
}

type ContainerCallAcceptedEvent struct{}

// implement fmt.Stringer interface
func (ccae ContainerCallAcceptedEvent) String() string {
	return "ContainerCallAccepted"
}

type OpCallAcceptedEvent struct{}

// implement fmt.Stringer interface
func (ocae OpCallAcceptedEvent) String() string {
	return "OpCallAccepted"
}

type ParallelCallAcceptedEvent struct{}

// implement fmt.Stringer interface
func (pcae ParallelCallAcceptedEvent) String() string {
	return "ParallelCallAccepted"
}

type SerialCallAcceptedEvent struct{}

// implement fmt.Stringer interface
func (scae SerialCallAcceptedEvent) String() string {
	return "SerialCallAccepted"
}
