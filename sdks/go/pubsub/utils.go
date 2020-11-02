package pubsub

import (
	"github.com/opctl/opctl/sdks/go/model"
)

func isRootCallIDExcludedByFilter(
	rootCallID string,
	filter model.EventFilter,
) bool {
	if nil != filter.Roots {
		isMatchFound := false
		for _, includedRootCallID := range filter.Roots {
			if includedRootCallID == rootCallID {
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

func getEventRootCallID(
	event model.Event,
) string {
	switch {
	case nil != event.CallEnded:
		return event.CallEnded.RootCallID
	case nil != event.ContainerStdErrWrittenTo:
		return event.ContainerStdErrWrittenTo.RootCallID
	case nil != event.ContainerStdOutWrittenTo:
		return event.ContainerStdOutWrittenTo.RootCallID
	case nil != event.OpKillRequested:
		return event.OpKillRequested.Request.RootCallID
	case nil != event.CallStarted:
		return event.CallStarted.RootCallID
	default:
		// use empty guid for unknown events
		return "00000000-0000-0000-0000-000000000000"
	}
}
