package core

import (
	"fmt"
	"github.com/opspec-io/engine/util/colorer"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"os"
	"time"
)

// allows mocking output
type output interface {
	Event(event *model.Event)
}

func newOutput(
	colorer colorer.Colorer,
) output {
	return &_output{
		colorer: colorer,
	}
}

type _output struct {
	colorer colorer.Colorer
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
	fmt.Print(
		this.colorer.Info(
			"ContainerExited Id='%v' OpRef='%v' ExitCode='%v' Timestamp='%v'\n",
			event.ContainerExited.ContainerId,
			event.ContainerExited.OpRef,
			event.ContainerExited.ExitCode,
			event.Timestamp.Format(time.RFC3339),
		),
	)
}

func (this _output) containerStarted(event *model.Event) {
	fmt.Print(
		this.colorer.Info(
			"ContainerStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
			event.ContainerStarted.ContainerId,
			event.ContainerStarted.OpRef,
			event.Timestamp.Format(time.RFC3339),
		),
	)
}

func (this _output) containerStdErrWrittenTo(event *model.Event) {
	fmt.Fprintf(os.Stderr, "%v \n", string(event.ContainerStdErrWrittenTo.Data))
}

func (this _output) containerStdOutWrittenTo(event *model.Event) {
	fmt.Fprintf(os.Stdout, "%v \n", string(event.ContainerStdOutWrittenTo.Data))
}

func (this _output) opEncounteredError(event *model.Event) {
	fmt.Print(
		this.colorer.Error(
			"OpEncounteredError Id='%v' OpRef='%v' Timestamp='%v' Msg='%v'\n",
			event.OpEncounteredError.OpId,
			event.OpEncounteredError.OpRef,
			event.Timestamp.Format(time.RFC3339),
			event.OpEncounteredError.Msg,
		),
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
		fmt.Print(this.colorer.Success(message))
	case model.OpOutcomeKilled:
		fmt.Print(this.colorer.Info(message))
	default:
		fmt.Print(this.colorer.Error(message))
	}
}

func (this _output) opStarted(event *model.Event) {
	fmt.Print(
		this.colorer.Info(
			"OpStarted Id='%v' OpRef='%v' Timestamp='%v'\n",
			event.OpStarted.OpId,
			event.OpStarted.OpRef,
			event.Timestamp.Format(time.RFC3339),
		),
	)
}
