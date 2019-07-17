import { eventStreamGet as opctlEventStreamGet } from '@opctl/sdk/src/api/client'
import uuidV4 from 'uuid/v4'
import isFiltered from './isFiltered'
import ModelEvent from '@opctl/sdk/src/model/event'
import { EventFilter } from '@opctl/sdk/src/api/client/events/stream'
import { toast } from 'react-toastify'

interface Subscription {
  filter: EventFilter
  onEvent: (event: ModelEvent) => any
}

const events = [] as ModelEvent[]
const subscriptions = {} as { [subscriptionId: string]: Subscription }

export default class EventStore {
  constructor(
    apiBaseUrl: string
  ) {
    opctlEventStreamGet(
      apiBaseUrl,
      {
        onEvent: (event: ModelEvent) => this.add(event),
        onError: (error: Error) => toast.error(error.message)
      }
    )
  }

  add(
    event: ModelEvent
  ) {
    events.push(event)
    Object.values(subscriptions).forEach(subscription => {
      if (subscription.onEvent) {
        if (!isFiltered(subscription.filter, event)) {
          subscription.onEvent(event)
        }
      }
    })
  }

  /**
   * subscribes to events
   * @param onEvent
   * @param {{roots,since}} filter
   * @param {Function} onEvent
   * @returns {function()} cancel; cancels any further calls to onEvent
   */
  getStream(
    filter: EventFilter,
    onEvent: (event: ModelEvent) => any
  ): () => void {
    const subscriptionId = uuidV4()

    events.forEach(event => {
      if (onEvent) {
        if (!isFiltered(filter, event)) {
          onEvent(event)
        }
      }
    })
    subscriptions[subscriptionId] = { onEvent, filter }

    return () => delete subscriptions[subscriptionId]
  }
}