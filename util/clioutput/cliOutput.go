package clioutput

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliOutput

import (
	"fmt"
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"io"
	"time"
)

// allows mocking/faking output
type CliOutput interface {
	// outputs a msg requiring attention
	Attention(format string, values ...interface{})

	// outputs an error msg
	Error(format string, values ...interface{})

	// outputs an event
	// @TODO: not generic
	Event(event *model.Event)

	// outputs an info msg
	Info(format string, values ...interface{})

	// outputs a success msg
	Success(format string, values ...interface{})
}

func New(
	colorer colorer.Colorer,
	errWriter io.Writer,
	stdWriter io.Writer,
) CliOutput {
	return _cliOutput{
		colorer:   colorer,
		errWriter: errWriter,
		stdWriter: stdWriter,
	}
}

type _cliOutput struct {
	colorer   colorer.Colorer
	errWriter io.Writer
	stdWriter io.Writer
}

func (this _cliOutput) Attention(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.colorer.Attention(format, values...))
}

func (this _cliOutput) Error(format string, values ...interface{}) {
	fmt.Fprintln(this.errWriter, this.colorer.Error(format, values...))
}

func (this _cliOutput) Event(event *model.Event) {
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

func (this _cliOutput) containerExited(event *model.Event) {
	this.Info(
		"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
		event.ContainerExited.ContainerId,
		event.ContainerExited.OpRef,
		event.ContainerExited.ExitCode,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) containerStarted(event *model.Event) {
	this.Info(
		"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
		event.ContainerStarted.ContainerId,
		event.ContainerStarted.OpRef,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) containerStdErrWrittenTo(event *model.Event) {
	fmt.Fprintln(this.errWriter, string(event.ContainerStdErrWrittenTo.Data))
}

func (this _cliOutput) containerStdOutWrittenTo(event *model.Event) {
	fmt.Fprintln(this.stdWriter, string(event.ContainerStdOutWrittenTo.Data))
}

func (this _cliOutput) opEncounteredError(event *model.Event) {
	this.Error(
		"OpEncounteredError Id='%v' OpRef='%v' Timestamp='%v' Msg='%v'\n",
		event.OpEncounteredError.OpId,
		event.OpEncounteredError.OpRef,
		event.Timestamp.Format(time.RFC3339),
		event.OpEncounteredError.Msg,
	)
}

func (this _cliOutput) opEnded(event *model.Event) {
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

func (this _cliOutput) opStarted(event *model.Event) {
	this.Info(
		"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
		event.OpStarted.OpId,
		event.OpStarted.OpRef,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) Info(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.colorer.Info(format, values...))
}

func (this _cliOutput) Success(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.colorer.Success(format, values...))
}
