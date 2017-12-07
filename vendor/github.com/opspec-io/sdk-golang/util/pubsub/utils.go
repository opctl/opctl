package pubsub

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
)

func isRootOpIdExcludedByFilter(
	rootOpId string,
	filter *model.EventFilter,
) bool {
	if nil != filter && nil != filter.Roots {
		isMatchFound := false
		for _, includedRootOpId := range filter.Roots {
			if includedRootOpId == rootOpId {
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

func getEventRootOpId(
	event *model.Event,
) string {
	switch {
	case nil != event.ContainerExited:
		return event.ContainerExited.RootOpId
	case nil != event.ContainerStarted:
		return event.ContainerStarted.RootOpId
	case nil != event.ContainerStdErrWrittenTo:
		return event.ContainerStdErrWrittenTo.RootOpId
	case nil != event.ContainerStdOutWrittenTo:
		return event.ContainerStdOutWrittenTo.RootOpId
	case nil != event.OpErred:
		return event.OpErred.RootOpId
	case nil != event.OpEnded:
		return event.OpEnded.RootOpId
	case nil != event.OpStarted:
		return event.OpStarted.RootOpId
	case nil != event.ParallelCallEnded:
		return event.ParallelCallEnded.RootOpId
	case nil != event.SerialCallEnded:
		return event.SerialCallEnded.RootOpId
	default:
		panic(fmt.Sprintf("Received unexpected event %v\n", event))
	}
}
