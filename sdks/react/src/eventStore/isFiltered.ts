import Event from '@opctl/sdk/src/model/event'
import { EventFilter } from '@opctl/sdk/src/api/client/events/stream'

const getEventRootOpId = (event: Event) => {
    if (event.callEnded) {
      return event.callEnded.rootOpId
    } else if (event.callKilled) {
      return event.callKilled.rootOpId
    } else if (event.containerExited) {
      return event.containerExited.rootOpId
    } else if (event.containerStarted) {
      return event.containerStarted.rootOpId
    } else if (event.containerStdErrWrittenTo) {
      return event.containerStdErrWrittenTo.rootOpId
    } else if (event.containerStdOutWrittenTo) {
      return event.containerStdOutWrittenTo.rootOpId
    } else if (event.opErred) {
      return event.opErred.rootOpId
    } else if (event.opEnded) {
      return event.opEnded.rootOpId
    } else if (event.opStarted) {
      return event.opStarted.rootOpId
    } else if (event.parallelCallEnded) {
      return event.parallelCallEnded.rootOpId
    } else if (event.serialCallEnded) {
      return event.serialCallEnded.rootOpId
    } else {
      throw new Error(`received unexpected event ${JSON.stringify(event)}`)
    }
  }

  export default function isFiltered (
      filter: EventFilter,
       event: Event
    ) {
    if (filter && Array.isArray(filter.roots)) {
      return !filter.roots.find(rootOpId => rootOpId === getEventRootOpId(event))
    }
    // @TODO: apply filter.since
    return false
  }