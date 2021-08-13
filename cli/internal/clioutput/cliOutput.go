package clioutput

import (
	"fmt"
	"io"
	"strings"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/sdks/go/model"
)

//CliOutput allows mocking/faking output
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
	case event.CallEnded != nil &&
		event.CallEnded.Call.Container != nil:
		this.containerExited(event)

	case event.CallStarted != nil &&
		event.CallStarted.Call.Container != nil:
		this.containerStarted(event)

	case event.ContainerStdErrWrittenTo != nil:
		this.containerStdErrWrittenTo(event.ContainerStdErrWrittenTo)

	case event.ContainerStdOutWrittenTo != nil:
		this.containerStdOutWrittenTo(event.ContainerStdOutWrittenTo)

	case event.CallEnded != nil &&
		event.CallEnded.Call.Op != nil:
		this.opEnded(event)

	case event.CallStarted != nil && event.CallStarted.Call.Op != nil:
		this.opStarted(event.CallStarted)
	}
}

func (this _cliOutput) containerExited(event *model.Event) {
	var color func(s string) string
	var writer io.Writer
	var message string
	switch event.CallEnded.Outcome {
	case model.OpOutcomeSucceeded:
		message = "exited"
		color = this.cliColorer.Success
		writer = this.stdWriter
	case model.OpOutcomeKilled:
		message = "killed"
		color = this.cliColorer.Info
		writer = this.stdWriter
	default:
		message = "crashed"
		color = this.cliColorer.Error
		writer = this.errWriter
	}

	if event.CallEnded.Call.Container.Image.Ref != nil {
		message = fmt.Sprintf("%s ", *event.CallEnded.Call.Container.Image.Ref) + message
	} else {
		message += "unknown container " + message
	}
	message = color(message)

	io.WriteString(
		writer,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Ref),
			message,
		),
	)
}

func (this _cliOutput) containerStarted(event *model.Event) {
	message := "started "
	if event.CallStarted.Call.Container.Image.Ref != nil {
		message += *event.CallStarted.Call.Container.Image.Ref
	} else {
		message += "unknown container"
	}

	io.WriteString(
		this.stdWriter,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.CallStarted.Call.ID, event.CallStarted.Ref),
			this.cliColorer.Info(message),
		),
	)
}

func (this _cliOutput) outputPrefix(id, opRef string) string {
	parts := []string{
		fmt.Sprintf("%.8s", fmt.Sprintf("%-8s", id)),
	}
	opRef = FormatOpRef(opRef)
	if opRef != "" {
		parts = append(parts, opRef)
	}
	return this.cliColorer.Muted("["+strings.Join(parts, " ")+"]") + " "
}

func (this _cliOutput) containerStdErrWrittenTo(event *model.ContainerStdErrWrittenTo) {
	io.WriteString(
		this.errWriter,
		fmt.Sprintf(
			"%s%s",
			this.outputPrefix(event.ContainerID, event.OpRef),
			event.Data,
		),
	)
}

func (this _cliOutput) containerStdOutWrittenTo(event *model.ContainerStdOutWrittenTo) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintf(
			"%s%s",
			this.outputPrefix(event.ContainerID, event.OpRef),
			event.Data,
		),
	)
}

func (this _cliOutput) opEnded(event *model.Event) {
	var color func(s string) string
	var writer io.Writer
	var message string
	switch event.CallEnded.Outcome {
	case model.OpOutcomeSucceeded:
		message = "succeeded"
		color = this.cliColorer.Success
		writer = this.stdWriter
	case model.OpOutcomeKilled:
		message = "killed"
		color = this.cliColorer.Info
		writer = this.stdWriter
	default:
		message = "failed"
		color = this.cliColorer.Error
		writer = this.errWriter
	}

	message = color(fmt.Sprintf("op %s", message))
	if event.CallEnded.Error != nil {
		message += color(":") + " " + event.CallEnded.Error.Message
	}

	io.WriteString(
		writer,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Call.Op.OpPath),
			message,
		),
	)
}

func (this _cliOutput) opStarted(event *model.CallStarted) {
	io.WriteString(
		this.stdWriter,
		fmt.Sprintf(
			"%s%s\n",
			this.outputPrefix(event.Call.ID, event.Call.Op.OpPath),
			this.cliColorer.Info("started op"),
		),
	)
}

func (this _cliOutput) info(s string) {
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
