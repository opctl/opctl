package pubsub

import (
	"fmt"
	"github.com/opctl/sdk-golang/model"
)

func isRootOpIDExcludedByFilter(
	rootOpID string,
	filter model.EventFilter,
) bool {
	if nil != filter.Roots {
		isMatchFound := false
		for _, includedRootOpID := range filter.Roots {
			if includedRootOpID == rootOpID {
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

func getEventRootOpID(
	event model.Event,
) string {
	switch {
	case nil != event.CallEnded:
		return event.CallEnded.RootCallID
	case nil != event.ContainerExited:
		return event.ContainerExited.RootOpID
	case nil != event.ContainerStarted:
		return event.ContainerStarted.RootOpID
	case nil != event.ContainerStdErrWrittenTo:
		return event.ContainerStdErrWrittenTo.RootOpID
	case nil != event.ContainerStdOutWrittenTo:
		return event.ContainerStdOutWrittenTo.RootOpID
	case nil != event.OpErred:
		return event.OpErred.RootOpID
	case nil != event.OpEnded:
		return event.OpEnded.RootOpID
	case nil != event.OpStarted:
		return event.OpStarted.RootOpID
	case nil != event.ParallelCallEnded:
		return event.ParallelCallEnded.RootOpID
	case nil != event.SerialCallEnded:
		return event.SerialCallEnded.RootOpID
	default:
		panic(fmt.Sprintf("Received unexpected event %v\n", event))
	}
}
