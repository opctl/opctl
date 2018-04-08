const WebSocket = require('universal-websocket-client')

class WebSocketClient {
  /**
   * Initialize a WebSocket client which reconnects when onclose returns unexpected code
   *
   * @param {String} address The URL to which to connect
   */
  constructor (address) {
    this.address = address
    this.client = {}
    this._initClient()
  }

  _initClient () {
    const prevClient = this.client
    const nextClient = new WebSocket(this.address)

    nextClient.onclose = prevClient.onclose
    nextClient.onerror = prevClient.onerror
    nextClient.onmessage = prevClient.onmessage
    nextClient.onopen = prevClient.onopen

    this.client = nextClient
  }

  close () {
    this.isClosed = true
    this.client.close(1000)
  }

  set onclose (callback) {
    this.client.onclose = closeEvent => {
      if (!this.isClosed && closeEvent.code !== 1000) {
        setTimeout(() => this._initClient(), 5000)
      }
      callback(closeEvent)
    }
  }

  set onerror (callback) {
    this.client.onerror = callback
  }

  set onmessage (callback) {
    this.client.onmessage = callback
  }

  set onopen (callback) {
    this.client.onopen = callback
  }

  send (data) {
    this.client.send(data)
  }
}

module.exports = WebSocketClient
