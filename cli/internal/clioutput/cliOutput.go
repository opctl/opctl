package clioutput

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/sdks/go/model"
)

// CliOutput allows mocking/faking output
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

// NewInteractive prefixes each line of output from background containers.
func NewInteractive(
	cliColorer clicolorer.CliColorer,
	errWriter, stdWriter io.Writer,
	writeFocused func(*model.Event) bool,
) CliOutput {
	return _cliOutput{
		cliColorer:   cliColorer,
		errWriter:    errWriter,
		stdWriter:    stdWriter,
		stdPartial:   make(map[string]bool),
		errPartial:   make(map[string]bool),
		writeFocused: writeFocused,
	}
}

type _cliOutput struct {
	cliColorer             clicolorer.CliColorer
	errWriter              io.Writer
	stdWriter              io.Writer
	stdPartial, errPartial map[string]bool
	writeFocused           func(*model.Event) bool
}

func (clio _cliOutput) DisableColor() {
	clio.cliColorer.DisableColor()
}

func (clio _cliOutput) Attention(s string) {
	io.WriteString(
		clio.stdWriter,
		fmt.Sprintln(
			clio.cliColorer.Attention(s),
		),
	)
}

func (clio _cliOutput) Warning(s string) {
	io.WriteString(
		clio.stdWriter,
		fmt.Sprintln(
			clio.cliColorer.Error(s),
		),
	)
}

func (clio _cliOutput) Error(s string) {
	io.WriteString(
		clio.errWriter,
		fmt.Sprintln(
			clio.cliColorer.Error(s),
		),
	)
}

func (clio _cliOutput) Event(event *model.Event) {
	switch {
	case event.CallEnded != nil &&
		event.CallEnded.Call.Container != nil:
		id := event.CallEnded.Call.Container.ContainerID
		delete(clio.stdPartial, id)
		delete(clio.errPartial, id)
		clio.containerExited(event)

	case event.CallStarted != nil &&
		event.CallStarted.Call.Container != nil:
		clio.containerStarted(event)

	case event.ContainerStdErrWrittenTo != nil:
		out := event.ContainerStdErrWrittenTo
		clio.containerOutput(event, clio.errWriter, clio.errPartial, out.ContainerID, out.OpRef, out.Data)

	case event.ContainerStdOutWrittenTo != nil:
		out := event.ContainerStdOutWrittenTo
		clio.containerOutput(event, clio.stdWriter, clio.stdPartial, out.ContainerID, out.OpRef, out.Data)

	case event.CallEnded != nil && event.CallEnded.Call.Op != nil:
		clio.opEnded(event)

	case event.CallStarted != nil && event.CallStarted.Call.Op != nil:
		clio.opStarted(event.CallStarted)

	case event.CallEnded != nil && event.CallEnded.Error != nil:
		clio.callEnded(event)
	}
}

func (clio _cliOutput) containerExited(event *model.Event) {
	var color func(s string) string
	var writer io.Writer
	var message string
	switch event.CallEnded.Outcome {
	case model.OpOutcomeSucceeded:
		message = "exited"
		color = clio.cliColorer.Success
		writer = clio.stdWriter
	case model.OpOutcomeKilled:
		message = "killed"
		color = clio.cliColorer.Info
		writer = clio.stdWriter
	default:
		message = "crashed"
		color = clio.cliColorer.Error
		writer = clio.errWriter
	}

	if event.CallEnded.Call.Container.Image.Ref != nil {
		message = fmt.Sprintf("%s ", *event.CallEnded.Call.Container.Image.Ref) + message
	} else {
		message += "unknown container " + message
	}
	message = color(message)
	if event.CallEnded.Error != nil {
		message += color(":") + " " + event.CallEnded.Error.Message
	}

	io.WriteString(
		writer,
		fmt.Sprintf(
			"%s%s\n",
			clio.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Ref),
			message,
		),
	)
}

func (clio _cliOutput) containerStarted(event *model.Event) {
	message := "started "
	if event.CallStarted.Call.Container.Image.Ref != nil {
		message += *event.CallStarted.Call.Container.Image.Ref
	} else {
		message += "unknown container"
	}

	io.WriteString(
		clio.stdWriter,
		fmt.Sprintf(
			"%s%s\n",
			clio.outputPrefix(event.CallStarted.Call.ID, event.CallStarted.Ref),
			clio.cliColorer.Info(message),
		),
	)
}

func (clio _cliOutput) outputPrefix(id, opRef string) string {
	parts := []string{
		fmt.Sprintf("%.8s", fmt.Sprintf("%-8s", id)),
	}
	opRef = FormatOpRef(opRef)
	if opRef != "" {
		parts = append(parts, opRef)
	}
	return clio.cliColorer.Muted("["+strings.Join(parts, " ")+"]") + " "
}

func (clio _cliOutput) callEnded(event *model.Event) {
	var color func(s string) string
	var writer io.Writer
	var message string
	switch event.CallEnded.Outcome {
	default:
		message = "failed"
		color = clio.cliColorer.Error
		writer = clio.errWriter
	}

	message = color(fmt.Sprintf("op %s", message))
	if event.CallEnded.Error != nil {
		message += color(":") + " " + event.CallEnded.Error.Message
	}

	io.WriteString(
		writer,
		fmt.Sprintf(
			"%s%s\n",
			clio.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Ref),
			message,
		),
	)
}

func (clio _cliOutput) containerOutput(
	event *model.Event,
	writer io.Writer,
	partial map[string]bool,
	id, opRef string,
	data []byte,
) {
	if partial == nil {
		fmt.Fprintf(writer, "%s%s", clio.outputPrefix(id, opRef), data)
		return
	}
	if len(data) == 0 {
		return
	}
	atLineStart := !partial[id]
	if data[len(data)-1] == '\n' {
		delete(partial, id)
	} else {
		partial[id] = true
	}
	if clio.writeFocused != nil && clio.writeFocused(event) {
		return
	}
	prefix := clio.outputPrefix(id, opRef)
	for line := range bytes.Lines(data) {
		if atLineStart {
			io.WriteString(writer, prefix)
		}
		writer.Write(line)
		atLineStart = true
	}
}

func (clio _cliOutput) opEnded(event *model.Event) {
	var color func(s string) string
	var writer io.Writer
	var message string
	switch event.CallEnded.Outcome {
	case model.OpOutcomeSucceeded:
		message = "succeeded"
		color = clio.cliColorer.Success
		writer = clio.stdWriter
	case model.OpOutcomeKilled:
		message = "killed"
		color = clio.cliColorer.Info
		writer = clio.stdWriter
	default:
		message = "failed"
		color = clio.cliColorer.Error
		writer = clio.errWriter
	}

	message = color(fmt.Sprintf("op %s", message))
	if event.CallEnded.Error != nil {
		message += color(":") + " " + event.CallEnded.Error.Message
	}

	io.WriteString(
		writer,
		fmt.Sprintf(
			"%s%s\n",
			clio.outputPrefix(event.CallEnded.Call.ID, event.CallEnded.Call.Op.OpPath),
			message,
		),
	)
}

func (clio _cliOutput) opStarted(event *model.CallStarted) {
	io.WriteString(
		clio.stdWriter,
		fmt.Sprintf(
			"%s%s\n",
			clio.outputPrefix(event.Call.ID, event.Call.Op.OpPath),
			clio.cliColorer.Info("started op"),
		),
	)
}

func (clio _cliOutput) Success(s string) {
	io.WriteString(
		clio.stdWriter,
		fmt.Sprintln(
			clio.cliColorer.Success(s),
		),
	)
}
