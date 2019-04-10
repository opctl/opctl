import _Websocket from '../_Websocket'
import Event from '../../../model/event'

export interface EventFilter {
    roots: string[]
}

export interface Options {
    filter: EventFilter

    /**
     * Callback invoked when event stream closes
     */
    onClose: () => void

    /**
     * Callback invoked when an error is received
     */
    onError: () => void

    /**
     * Callback invoked for each event received
     */
    onEvent: (event: Event) => void

    /**
     * Callback invoked when event stream opens
     */
    onOpen: () => void
}

/**
 * Gets an event stream
 * implements https://github.com/opctl/spec/blob/0.1.5/spec/node-api.spec.yml#L7
 * interface loosely based on: https://html.spec.whatwg.org/multipage/web-sockets.html#network
 *
 * @return {Function} stream closer
 */
export default function stream(
    apiBaseUrl: string,
    {
        filter,
        onClose = () => { },
        onError = () => { },
        onEvent,
        onOpen = () => { }
    }: Options
): () => void {
    // construct request URL
    let queryParts = []
    const defaultedFilter = Object.assign(
        { roots: [] },
        filter
    )
    queryParts.push(`roots=${defaultedFilter.roots.map(root => encodeURIComponent(root)).join(',')}`)

    // enable backpressure
    queryParts.push(`ack`)

    // construct websocket client
    const baseWebsocketUrl = apiBaseUrl.replace(/^http/, 'ws')
    const webSocketClient = new _Websocket(`${baseWebsocketUrl}/events/stream?${queryParts.join('&')}`)
    webSocketClient.onclose = onClose
    webSocketClient.onerror = onError
    webSocketClient.onmessage = (msg: { data: string }) => {
        // we ack every message which serves as backpressure so the server doesn't flood us
        setTimeout(() => webSocketClient.send(''), 0)

        onEvent(JSON.parse(msg.data))
    }
    webSocketClient.onopen = onOpen

    return () => webSocketClient.close()
}