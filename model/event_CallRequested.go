package model

import (
	"fmt"
	"strings"
)

// CallRequestedEvent represents a request to carry out a call.
type CallRequestedEvent struct {
	*CallEventBase
	Container *ContainerCallRequestedEvent `json:"container,omitempty" yaml:"container,omitempty"`
	Op        *OpCallRequestedEvent        `json:"op,omitempty" yaml:"op,omitempty"`
	Parallel  *ParallelCallRequestedEvent  `json:"parallel,omitempty" yaml:"parallel,omitempty"`
	Serial    *SerialCallRequestedEvent    `json:"serial,omitempty" yaml:"serial,omitempty"`
}

// implement fmt.Stringer interface
func (cre CallRequestedEvent) String() string {
	var pre string
	switch {
	case nil != cre.Container:
		pre = cre.Container.String()
	case nil != cre.Op:
		pre = cre.Op.String()
	case nil != cre.Parallel:
		pre = cre.Parallel.String()
	case nil != cre.Serial:
		pre = cre.Serial.String()
	}

	return strings.Join(
		[]string{
			pre,
			cre.CallEventBase.String(),
		},
		" ",
	)
}

type ContainerCallRequestedEvent struct {
	ImageRef string
}

// implement fmt.Stringer interface
func (ccre ContainerCallRequestedEvent) String() string {
	return strings.Join(
		[]string{
			"ContainerCallRequested",
			fmt.Sprintf("ImageRef='%v'", ccre.ImageRef),
		},
		" ",
	)
}

type OpCallRequestedEvent struct {
	PkgRef string
}

// implement fmt.Stringer interface
func (ocre OpCallRequestedEvent) String() string {
	return strings.Join(
		[]string{
			"OpCallRequested",
			fmt.Sprintf("PkgRef='%v'", ocre.PkgRef),
		},
		" ",
	)
}

type ParallelCallRequestedEvent struct{}

// implement fmt.Stringer interface
func (pcre ParallelCallRequestedEvent) String() string {
	return "ParallelCallRequested"
}

type SerialCallRequestedEvent struct{}

// implement fmt.Stringer interface
func (scre SerialCallRequestedEvent) String() string {
	return "SerialCallRequested"
}
