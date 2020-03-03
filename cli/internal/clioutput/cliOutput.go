package clioutput

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"fmt"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/sdks/go/model"
	"io"
	"time"
)

//CliOutput allows mocking/faking output
//counterfeiter:generate -o fakes/cliOutput.go . CliOutput
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
	case nil != event.OpErred:
		this.opErred(event)
	case nil != event.OpEnded:
		this.opEnded(event)
	case nil != event.OpStarted:
		this.opStarted(event)
	}
}

func (this _cliOutput) containerExited(event *model.Event) {
	this.Info(
		"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
		event.ContainerExited.ContainerID,
		event.ContainerExited.OpRef,
		event.ContainerExited.ExitCode,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) containerStarted(event *model.Event) {
	this.Info(
		"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
		event.ContainerStarted.ContainerID,
		event.ContainerStarted.OpRef,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) containerStdErrWrittenTo(event *model.Event) {
	fmt.Fprint(this.errWriter, string(event.ContainerStdErrWrittenTo.Data))
}

func (this _cliOutput) containerStdOutWrittenTo(event *model.Event) {
	fmt.Fprint(this.stdWriter, string(event.ContainerStdOutWrittenTo.Data))
}

func (this _cliOutput) opErred(event *model.Event) {
	this.Error(
		"OpErred Id='%v' OpRef='%v' Timestamp='%v' Msg='%v'\n",
		event.OpErred.OpID,
		event.OpErred.OpRef,
		event.Timestamp.Format(time.RFC3339),
		event.OpErred.Msg,
	)
}

func (this _cliOutput) opEnded(event *model.Event) {
	message := fmt.Sprintf(
		"OpEnded Id='%v' OpRef='%v' Outcome='%v' Timestamp='%v'\n",
		event.OpEnded.OpID,
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
		event.OpStarted.OpID,
		event.OpStarted.OpRef,
		event.Timestamp.Format(time.RFC3339),
	)
}

func (this _cliOutput) Info(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.cliColorer.Info(format, values...))
}

func (this _cliOutput) Success(format string, values ...interface{}) {
	fmt.Fprintln(this.stdWriter, this.cliColorer.Success(format, values...))
}
