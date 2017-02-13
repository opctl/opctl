package core

//go:generate counterfeiter -o ./fakeOutput.go --fake-name fakeOutput ./ output

import (
	"fmt"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"io"
	"time"
)

// allows mocking/faking output
type output interface {
	// outputs a msg requiring attention
	Attention(format string, values ...interface{})

	// outputs an error msg
	Error(format string, values ...interface{})

	// outputs an event
	Event(event *model.Event)

	// outputs an info msg
	Info(format string, values ...interface{})

	// outputs a success msg
	Success(format string, values ...interface{})
}

func newOutput(
	colorer colorer.Colorer,
	errWriter io.Writer,
	stdWriter io.Writer,
) output {
	return &_output{
		colorer:   colorer,
		errWriter: errWriter,
		stdWriter: stdWriter,
	}
}

type _output struct {
	colorer   colorer.Colorer
	errWriter io.Writer
	stdWriter io.Writer
}

func (this _output) Attention(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.colorer.Attention(format, values...))
}

func (this _output) Error(format string, values ...interface{}) {
	fmt.Fprintln(this.errWriter, this.colorer.Error(format, values...))
}

func (this _output) Event(event *model.Event) {
	switch {
	case nil != event.ContainerExited:
		this.containerExited(event)
	case nil != event.ContainerStarted:
		this.containerStarted(event)
	case nil != event.ContainerStdErrWrittenTo:
		this.containerStdErrWrittenTo(event)
	case nil != event.ContainerStdOutWrittenTo:
		this.containerStdOutWrittenTo(event)
	case nil != event.OpEncounteredError:
		this.opEncounteredError(event)
	case nil != event.OpEnded:
		this.opEnded(event)
	case nil != event.OpStarted:
		this.opStarted(event)
	}
}

func (this _output) containerExited(event *model.Event) {
	this.Info(
		"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
		event.ContainerExited.ContainerId,
		event.ContainerExited.OpRef,
		event.ContainerExited.ExitCode,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _output) containerStarted(event *model.Event) {
	this.Info(
		"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
		event.ContainerStarted.ContainerId,
		event.ContainerStarted.OpRef,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _output) containerStdErrWrittenTo(event *model.Event) {
	fmt.Fprintln(this.errWriter, string(event.ContainerStdErrWrittenTo.Data))
}

func (this _output) containerStdOutWrittenTo(event *model.Event) {
	fmt.Fprintln(this.stdWriter, string(event.ContainerStdOutWrittenTo.Data))
}

func (this _output) opEncounteredError(event *model.Event) {
	this.Error(
		"OpEncounteredError Id='%v' OpRef='%v' Timestamp='%v' Msg='%v'\n",
		event.OpEncounteredError.OpId,
		event.OpEncounteredError.OpRef,
		event.Timestamp.Format(time.RFC3339),
		event.OpEncounteredError.Msg,
	)
}

func (this _output) opEnded(event *model.Event) {
	message := fmt.Sprintf(
		"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
		event.OpEnded.OpId,
		event.OpEnded.OpRef,
		event.OpEnded.Outcome,
		event.Timestamp.Format(time.RFC3339),
	)
	switch event.OpEnded.Outcome {
	case model.OpOutcomeSucceeded:
		this.Success(message)
	case model.OpOutcomeKilled:
		this.Info(message)
	default:
		this.Error(message)
	}
}

func (this _output) opStarted(event *model.Event) {
	this.Info(
		"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
		event.OpStarted.OpId,
		event.OpStarted.OpRef,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _output) Info(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.colorer.Info(format, values...))
}

func (this _output) Success(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.colorer.Success(format, values...))
}
