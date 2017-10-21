'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var fetch = require('node-fetch');
var WebSocket = require('universal-websocket-client');

var ApiClient = function () {
  function ApiClient(_ref) {
    var _ref$baseUrl = _ref.baseUrl,
        baseUrl = _ref$baseUrl === undefined ? 'http://localhost:42224' : _ref$baseUrl;

    _classCallCheck(this, ApiClient);

    this.baseHttpUrl = baseUrl;
    this.baseWebsocketUrl = this.baseHttpUrl.replace(/^http/, 'ws');
  }

  /**
   * Asserts response.status is in the range of successful status codes
   * @param response
   * @return {*}
   * @private
   */


  _createClass(ApiClient, [{
    key: '_assertStatusSuccessful',
    value: function _assertStatusSuccessful(response) {
      if (response.status >= 200 && response.status < 300) {
        return response;
      } else {
        return response.text().then(function (errorMsg) {
          var error = new Error(errorMsg);
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

  }, {
    key: 'event_stream_get',
    value: function event_stream_get(_ref2) {
      var filter = _ref2.filter,
          _ref2$onClose = _ref2.onClose,
          onClose = _ref2$onClose === undefined ? function (err) {} : _ref2$onClose,
          _ref2$onError = _ref2.onError,
          onError = _ref2$onError === undefined ? function (err) {} : _ref2$onError,
          onEvent = _ref2.onEvent,
          _ref2$onOpen = _ref2.onOpen,
          onOpen = _ref2$onOpen === undefined ? function () {} : _ref2$onOpen;

      // construct request URL
      var queryParts = [];
      var defaultedFilter = Object.assign({ roots: [] }, filter);
      queryParts.push('roots=' + defaultedFilter.roots.map(function (root) {
        return encodeURIComponent(root);
      }).join(','));

      // construct websocket client
      var webSocket = new WebSocket(this.baseWebsocketUrl + '/events/stream?' + queryParts.join('&'));
      webSocket.onclose = onClose;
      webSocket.onerror = onError;
      webSocket.onmessage = function (msg) {
        return onEvent(JSON.parse(msg.data));
      };
      webSocket.onopen = onOpen;

      return function () {
        return webSocket.close();
      };
    }

    /**
     * Gets liveness of node
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L60
     * @return {Promise.<fetch.Response>}
     */

  }, {
    key: 'liveness_get',
    value: function liveness_get() {
      return fetch(this.baseHttpUrl + '/liveness').then(this._assertStatusSuccessful).then(function (response) {
        return response.text();
      });
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

  }, {
    key: 'op_start',
    value: function op_start(opStartReq) {
      return fetch(this.baseHttpUrl + '/ops/starts', {
        method: 'POST',
        body: JSON.stringify(opStartReq)
      }).then(this._assertStatusSuccessful).then(function (response) {
        return response.text();
      });
    }

    /**
     * Kills an op
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L139
     * @param {Object} opKillReq
     * @param {Object} opKillReq.opId
     * @return {Promise.<null>}
     */

  }, {
    key: 'op_kill',
    value: function op_kill(opKillReq) {
      return fetch(this.baseHttpUrl + '/ops/kills', {
        method: 'POST',
        body: JSON.stringify(opKillReq)
      }).then(this._assertStatusSuccessful).then(function () {
        return null;
      });
    }

    /**
     * Gets pkg content at contentPath
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L242
     * @param pkgRef
     * @param contentPath
     * @return {Promise.<fetch.Response>}
     */

  }, {
    key: 'pkg_content_get',
    value: function pkg_content_get(_ref3) {
      var pkgRef = _ref3.pkgRef,
          contentPath = _ref3.contentPath;

      return fetch(this.baseHttpUrl + '/pkgs/' + encodeURIComponent(pkgRef) + '/contents/' + encodeURIComponent(contentPath)).then(this._assertStatusSuccessful);
    }

    /**
     * Lists pkg contents
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L178
     * @param pkgRef
     * @return {Promise.<Object[]>}}
     */

  }, {
    key: 'pkg_content_list',
    value: function pkg_content_list(_ref4) {
      var pkgRef = _ref4.pkgRef;

      return fetch(this.baseHttpUrl + '/pkgs/' + encodeURIComponent(pkgRef) + '/contents').then(this._assertStatusSuccessful).then(function (response) {
        return response.json();
      });
    }
  }]);

  return ApiClient;
}();

module.exports = ApiClient;