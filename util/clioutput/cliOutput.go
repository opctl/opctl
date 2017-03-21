package clioutput

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ CliOutput

import (
	"fmt"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opspec-io/sdk-golang/model"
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
	cliColorer clicolorer.CliColorer,
	errWriter io.Writer,
	stdWriter io.Writer,
) CliOutput {
	return _cliOutput{
		cliColorer: cliColorer,
		errWriter:  errWriter,
		stdWriter:  stdWriter,
	}
}

type _cliOutput struct {
	cliColorer clicolorer.CliColorer
	errWriter  io.Writer
	stdWriter  io.Writer
}

func (this _cliOutput) Attention(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.cliColorer.Attention(format, values...))
}

func (this _cliOutput) Error(format string, values ...interface{}) {
	fmt.Fprintln(this.errWriter, this.cliColorer.Error(format, values...))
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
		"ContainerExited Id='%v' PkgRef='%v' ExitCode='%v' Timestamp='%v'\n",
		event.ContainerExited.ContainerId,
		event.ContainerExited.PkgRef,
		event.ContainerExited.ExitCode,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) containerStarted(event *model.Event) {
	this.Info(
		"ContainerStarted Id='%v' PkgRef='%v' Timestamp='%v'\n",
		event.ContainerStarted.ContainerId,
		event.ContainerStarted.PkgRef,
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
		"OpEncounteredError Id='%v' PkgRef='%v' Timestamp='%v' Msg='%v'\n",
		event.OpEncounteredError.OpId,
		event.OpEncounteredError.PkgRef,
		event.Timestamp.Format(time.RFC3339),
		event.OpEncounteredError.Msg,
	)
}

func (this _cliOutput) opEnded(event *model.Event) {
	message := fmt.Sprintf(
		"OpEnded Id='%v' PkgRef='%v' Outcome='%v' Timestamp='%v'\n",
		event.OpEnded.OpId,
		event.OpEnded.PkgRef,
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
		"OpStarted Id='%v' PkgRef='%v' Timestamp='%v'\n",
		event.OpStarted.OpId,
		event.OpStarted.PkgRef,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) Info(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.cliColorer.Info(format, values...))
}

func (this _cliOutput) Success(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.cliColorer.Success(format, values...))
}
