package pubsub

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/pkg/model"
)

func isOgIdExcludedByFilter(
	ogid string,
	filter *model.EventFilter,
) bool {
	if nil != filter && nil != filter.OpGraphIds {
		isMatchFound := false
		for _, includedOgid := range filter.OpGraphIds {
			if includedOgid == ogid {
				isMatchFound = true
				break
			}
		}
		if !isMatchFound {
			return true
		}
	}
	return false
}

func getEventOpGraphId(
	event *model.Event,
) string {
	switch {
	case nil != event.ContainerExited:
		return event.ContainerExited.OpGraphId
	case nil != event.ContainerStarted:
		return event.ContainerStarted.OpGraphId
	case nil != event.ContainerStdErrWrittenTo:
		return event.ContainerStdErrWrittenTo.OpGraphId
	case nil != event.ContainerStdOutWrittenTo:
		return event.ContainerStdOutWrittenTo.OpGraphId
	case nil != event.OpEncounteredError:
		return event.OpEncounteredError.OpGraphId
	case nil != event.OpEnded:
		return event.OpEnded.OpGraphId
	case nil != event.OpStarted:
		return event.OpStarted.OpGraphId
	default:
		panic(fmt.Sprintf("Received unexpected event %v\n", event))
	}
}
