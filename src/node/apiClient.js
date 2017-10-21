const fetch = require('node-fetch');
const WebSocket = require('universal-websocket-client');

class ApiClient {
  constructor({baseUrl = 'http://localhost:42224'}) {
    this.baseHttpUrl = baseUrl;
    this.baseWebsocketUrl = this.baseHttpUrl.replace(/^http/, 'ws');
  }

  /**
   * Asserts response.status is in the range of successful status codes
   * @param response
   * @return {*}
   * @private
   */
  _assertStatusSuccessful(response) {
    if (response.status >= 200 && response.status < 300) {
      return response
    } else {
      return response.text().then(errorMsg => {
        const error = new Error(errorMsg);
        error.response = response;
        throw error;
      });
    }
  }

  /**
   * Invoked when event stream closes
   * @callback ApiClient~event_stream_get_onClose
   * @param {string} error
   */

  /**
   * Invoked when an error is received
   * @callback ApiClient~event_stream_get_onError
   * @param {string} error
   */

  /**
   * Invoked for each event received
   * @callback ApiClient~event_stream_get_onEvent
   * @param {Object} event https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L494
   */

  /**
   * Invoked when event stream opens
   * @callback ApiClient~event_stream_get_onOpen
   */

  /**
   * Gets an event stream
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L7
   * interface loosely based on: https://html.spec.whatwg.org/multipage/web-sockets.html#network
   *
   * @param {Object} [filter]
   * @param {Array<String>} [filter.roots]
   * @param {ApiClient~event_stream_get_onClose} onClose
   * @param {ApiClient~event_stream_get_onError} onError
   * @param {ApiClient~event_stream_get_onEvent} onEvent
   * @param {ApiClient~event_stream_get_onOpen} onOpen
   *
   * @return {Function} stream closer
   */
  event_stream_get({
                      filter,
                      onClose = err => {},
                      onError = err => {},
                      onEvent,
                      onOpen = () => {},
                    }) {
    // construct request URL
    let queryParts = [];
    const defaultedFilter = Object.assign({roots: []}, filter);
    queryParts.push(`roots=${defaultedFilter.roots.map(root => encodeURIComponent(root)).join(',')}`);

    // construct websocket client
    const webSocket = new WebSocket(`${this.baseWebsocketUrl}/events/stream?${queryParts.join('&')}`);
    webSocket.onclose = onClose;
    webSocket.onerror = onError;
    webSocket.onmessage = msg => onEvent(JSON.parse(msg.data));
    webSocket.onopen = onOpen;

    return () => webSocket.close();
  }

  /**
   * Gets liveness of node
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L60
   * @return {Promise.<fetch.Response>}
   */
  liveness_get() {
    return fetch(
      `${this.baseHttpUrl}/liveness`
    )
      .then(this._assertStatusSuccessful)
      .then(response => (response.text()));
  }

  /**
   * Starts an op
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L70
   * @param {Object} opStartReq
   * @param {Object} opStartReq.pkg
   * @param {String} opStartReq.pkg.ref
   * @param {Object} [opStartReq.pkg.pullCreds]
   * @param {String} [opStartReq.pkg.pullCreds.username]
   * @param {String} [opStartReq.pkg.pullCreds.password]
   * @param {Object} [opStartReq.args]
   * @return {Promise.<String>} id of the started op
   */
  op_start(opStartReq) {
    return fetch(`${this.baseHttpUrl}/ops/starts`, {
      method: 'POST',
      body: JSON.stringify(opStartReq),
    })
      .then(this._assertStatusSuccessful)
      .then(response => (response.text()))
  }

  /**
   * Kills an op
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L139
   * @param {Object} opKillReq
   * @param {Object} opKillReq.opId
   * @return {Promise.<null>}
   */
  op_kill(opKillReq) {
    return fetch(`${this.baseHttpUrl}/ops/kills`, {
      method: 'POST',
      body: JSON.stringify(opKillReq),
    })
      .then(this._assertStatusSuccessful)
      .then(() => null);
  }

  /**
   * Gets pkg content at contentPath
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L242
   * @param pkgRef
   * @param contentPath
   * @return {Promise.<fetch.Response>}
   */
  pkg_content_get({pkgRef, contentPath}) {
    return fetch(
      `${this.baseHttpUrl}/pkgs/${encodeURIComponent(pkgRef)}/contents/${encodeURIComponent(contentPath)}`
    )
      .then(this._assertStatusSuccessful);
  }

  /**
   * Lists pkg contents
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L178
   * @param pkgRef
   * @return {Promise.<Object[]>}}
   */
  pkg_content_list({pkgRef}) {
    return fetch(
      `${this.baseHttpUrl}/pkgs/${encodeURIComponent(pkgRef)}/contents`
    )
      .then(this._assertStatusSuccessful)
      .then(response => (response.json()));
  }
}

module.exports = ApiClient;
