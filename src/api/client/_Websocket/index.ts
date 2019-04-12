import UniversalWebsocketClient from 'universal-websocket-client'

export default interface CloseEvent {
  code: number
  reason: string
  wasClean: boolean
}

export default class _Websocket {
  /**
   * Initialize a WebSocket client which reconnects when onclose returns unexpected code
   *
   * @param {String} address The URL to which to connect
   */
  constructor(address: string) {
    this.address = address
    this.client = new UniversalWebsocketClient(this.address)
    this._initClient()
  }

  address: string
  client: any
  isClosed: boolean = false

  _initClient() {
    const prevClient = this.client
    const nextClient = new UniversalWebsocketClient(this.address)

    nextClient.onclose = prevClient.onclose
    nextClient.onerror = prevClient.onerror
    nextClient.onmessage = prevClient.onmessage
    nextClient.onopen = prevClient.onopen

    this.client = nextClient
  }

  close() {
    this.isClosed = true
    this.client.close(1000)
  }

  set onclose(
    callback: (closeEvent: CloseEvent) => void
  ) {
    this.client.onclose = (closeEvent: CloseEvent) => {
      if (!this.isClosed && closeEvent.code !== 1000) {
        setTimeout(() => this._initClient(), 5000)
      }
      callback(closeEvent)
    }
  }

  set onerror(
    callback: (error: Error) => void
  ) {
    this.client.onerror = callback
  }

  set onmessage(
    callback: (message: any) => void
  ) {
    this.client.onmessage = callback
  }

  set onopen(
    callback: () => void
  ) {
    this.client.onopen = callback
  }

  send(
    data: string
  ) {
    this.client.send(data)
  }
}