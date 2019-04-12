import getApiBaseUrl from '../getApiBaseUrl'
import { eventStreamGet } from '@opctl/sdk/lib/api/client'
import uuidV4 from 'uuid/v4'
import filterChecker from './filterChecker'

const events = [] as any[]
const subscriptions = [] as any[]

class EventStore {
  constructor () {
    eventStreamGet(
      getApiBaseUrl(),
      {
        onEvent: event => this.add(event)
      } as any
    )
  }

  add (event) {
    events.push(event)
    Object.values(subscriptions).forEach(subscription => {
      if (subscription.onEvent) {
        if (!filterChecker.isFiltered(subscription.filter, event)) {
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
  getStream ({
    filter,
    onEvent
  }) {
    const subscriptionId = uuidV4()

    events.forEach(event => {
      if (onEvent) {
        if (!filterChecker.isFiltered(filter, event)) {
          onEvent(event)
        }
      }
    })
    subscriptions[subscriptionId] = { onEvent, filter }

    return () => delete subscriptions[subscriptionId]
  }
}

export default new EventStore()
