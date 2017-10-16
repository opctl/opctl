'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var fetch = require('node-fetch');

var ApiClient = function () {
  function ApiClient(_ref) {
    var _ref$baseUrl = _ref.baseUrl,
        baseUrl = _ref$baseUrl === undefined ? 'http://localhost:42224' : _ref$baseUrl;

    _classCallCheck(this, ApiClient);

    this.baseUrl = baseUrl;
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
     * Gets liveness of node
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L60
     * @return {Promise.<fetch.Response>}
     */

  }, {
    key: 'liveness_get',
    value: function liveness_get() {
      return fetch(this.baseUrl + '/liveness').then(this._assertStatusSuccessful).then(function (response) {
        return response.text();
      });
    }

    /**
     * Starts an op
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L70
     * @param opStartReq
     * @return {Promise.<String>} id of the started op
     */

  }, {
    key: 'op_start',
    value: function op_start(opStartReq) {
      return fetch(this.baseUrl + '/ops/starts', {
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
     * @param opKillReq
     * @return {Promise.<null>}
     */

  }, {
    key: 'op_kill',
    value: function op_kill(opKillReq) {
      return fetch(this.baseUrl + '/ops/kills', {
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
    value: function pkg_content_get(_ref2) {
      var pkgRef = _ref2.pkgRef,
          contentPath = _ref2.contentPath;

      return fetch(this.baseUrl + '/pkgs/' + encodeURIComponent(pkgRef) + '/contents/' + encodeURIComponent(contentPath)).then(this._assertStatusSuccessful);
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
    value: function pkg_content_list(_ref3) {
      var pkgRef = _ref3.pkgRef;

      return fetch(this.baseUrl + '/pkgs/' + encodeURIComponent(pkgRef) + '/contents').then(this._assertStatusSuccessful).then(function (response) {
        return response.json();
      });
    }
  }]);

  return ApiClient;
}();

module.exports = ApiClient;