package pubsub

import (
	"github.com/opctl/opctl/sdks/go/model"
)

func isRootCallIDExcludedByFilter(
	rootCallID string,
	filter model.EventFilter,
) bool {
	if filter.Roots == nil {
		return false
	}

	for _, includedRootCallID := range filter.Roots {
		if includedRootCallID == rootCallID {
			return false
		}
	}

	return true
}

func getEventRootCallID(
	event model.Event,
) string {
	switch {
	case event.CallEnded != nil:
		return event.CallEnded.Call.RootID
	case event.ContainerStdErrWrittenTo != nil:
		return event.ContainerStdErrWrittenTo.RootCallID
	case event.ContainerStdOutWrittenTo != nil:
		return event.ContainerStdOutWrittenTo.RootCallID
	case event.CallKillRequested != nil:
		return event.CallKillRequested.Request.RootCallID
	case event.CallStarted != nil:
		return event.CallStarted.Call.RootID
	default:
		// use empty guid for unknown events
		return "00000000-0000-0000-0000-000000000000"
	}
}
