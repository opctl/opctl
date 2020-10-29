package clioutput

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"fmt"
	"io"
	"time"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/sdks/go/model"
)

//CliOutput allows mocking/faking output
//counterfeiter:generate -o fakes/cliOutput.go . CliOutput
type CliOutput interface {
	// outputs a msg requiring attention
	Attention(s string)

	// outputs an error msg
	Error(s string)

	// outputs an event
	// @TODO: not generic
	Event(event *model.Event)

	// outputs an info msg
	Info(s string)

	// outputs a success msg
	Success(s string)
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

func (this _cliOutput) Attention(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Attention(s),
		),
	)
}

func (this _cliOutput) Error(s string) {
	io.WriteString(
		this.errWriter,
		fmt.Sprintln(
			this.cliColorer.Error(s),
		),
	)
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
	case nil != event.CallEnded && event.CallEnded.CallType == model.CallTypeOp:
		this.opEnded(event)
	case nil != event.CallStarted && event.CallStarted.CallType == model.CallTypeOp:
		this.opStarted(event)
	}
}

func (this _cliOutput) containerExited(event *model.Event) {
	this.Info(
		fmt.Sprintf(
			"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
			event.ContainerExited.ContainerID,
			event.ContainerExited.OpRef,
			event.ContainerExited.ExitCode,
			event.Timestamp.Format(time.RFC3339),
		),
	)
}

func (this _cliOutput) containerStarted(event *model.Event) {
	this.Info(
		fmt.Sprintf(
			"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
			event.ContainerStarted.ContainerID,
			event.ContainerStarted.OpRef,
			event.Timestamp.Format(time.RFC3339),
		),
	)
}

func (this _cliOutput) containerStdErrWrittenTo(event *model.Event) {
	io.WriteString(this.errWriter, string(event.ContainerStdErrWrittenTo.Data))
}

func (this _cliOutput) containerStdOutWrittenTo(event *model.Event) {
	io.WriteString(this.stdWriter, string(event.ContainerStdOutWrittenTo.Data))
}

func (this _cliOutput) opEnded(event *model.Event) {
	err := ""
	if nil != event.CallEnded.Error {
		err = fmt.Sprintf(" Error='%v'", event.CallEnded.Error.Message)
	}
	message := fmt.Sprintf(
		"OpEnded Id='%v' OpRef='%v' Outcome='%v'%v Timestamp='%v'\n",
		event.CallEnded.CallID,
		event.CallEnded.Ref,
		event.CallEnded.Outcome,
		err,
		event.Timestamp.Format(time.RFC3339),
	)
	switch event.CallEnded.Outcome {
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
		fmt.Sprintf(
			"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
			event.CallStarted.CallID,
			event.CallStarted.OpRef,
			event.Timestamp.Format(time.RFC3339),
		),
	)
}

func (this _cliOutput) Info(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Info(s),
		),
	)
}

func (this _cliOutput) Success(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Success(s),
		),
	)
}
