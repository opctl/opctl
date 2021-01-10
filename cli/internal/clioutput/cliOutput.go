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
	// silently disables coloring
	DisableColor()

	// outputs a msg requiring attention
	Attention(s string)

	// outputs a warning message (looks like an error but on stdout)
	Warning(s string)

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

func (this _cliOutput) DisableColor() {
	this.cliColorer.DisableColor()
}

func (this _cliOutput) Attention(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Attention(s),
		),
	)
}

func (this _cliOutput) Warning(s string) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintln(
			this.cliColorer.Error(s),
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
	case nil != event.CallEnded &&
		nil == event.CallEnded.Call.Op &&
		nil == event.CallEnded.Call.Container &&
		nil != event.CallEnded.Error:
		this.error(event)

	case nil != event.CallEnded &&
		nil != event.CallEnded.Call.Container:
		this.containerExited(event)

	case nil != event.CallStarted &&
		nil != event.CallStarted.Call.Container:
		this.containerStarted(event)

	case nil != event.ContainerStdErrWrittenTo:
		this.containerStdErrWrittenTo(event)

	case nil != event.ContainerStdOutWrittenTo:
		this.containerStdOutWrittenTo(event)

	case nil != event.CallEnded &&
		nil != event.CallEnded.Call.Op:
		this.opEnded(event)

	case nil != event.CallStarted &&
		nil != event.CallStarted.Call.Op:
		this.opStarted(event)
	}
}

func (this _cliOutput) error(event *model.Event) {
	this.Error(
		fmt.Sprintf(
			"Error='%v' Id='%v' OpRef='%v' Timestamp='%v'\n",
			event.CallEnded.Error.Message,
			event.CallEnded.Call.ID,
			event.CallEnded.Ref,
			event.Timestamp.Format(time.RFC3339),
		),
	)
}

func (this _cliOutput) containerExited(event *model.Event) {
	err := ""
	if nil != event.CallEnded.Error {
		err = fmt.Sprintf(" Error='%v'", event.CallEnded.Error.Message)
	}

	imageRef := ""
	if nil != event.CallEnded.Call.Container.Image.Ref {
		imageRef = fmt.Sprintf(" ImageRef='%v'", *event.CallEnded.Call.Container.Image.Ref)
	}

	message := fmt.Sprintf(
		"ContainerExited Id='%v'%v Outcome='%v'%v Timestamp='%v'\n",
		event.CallEnded.Call.ID,
		imageRef,
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

func (this _cliOutput) containerStarted(event *model.Event) {
	imageRef := ""
	if nil != event.CallStarted.Call.Container.Image.Ref {
		imageRef = fmt.Sprintf(" ImageRef='%v'", *event.CallStarted.Call.Container.Image.Ref)
	}

	this.Info(
		fmt.Sprintf(
			"ContainerStarted Id='%v' OpRef='%v'%v Timestamp='%v'\n",
			event.CallStarted.Call.ID,
			event.CallStarted.Ref,
			imageRef,
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
		event.CallEnded.Call.ID,
		event.CallEnded.Call.Op.OpPath,
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
			event.CallStarted.Call.ID,
			event.CallStarted.Call.Op.OpPath,
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
