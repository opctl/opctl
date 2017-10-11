'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { step("next", value); }, function (err) { step("throw", err); }); } } return step("next"); }); }; }

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
   * Starts an op
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L70
   * @param opStartReq
   * @return {Promise.<String>} id of the started op
   */


  _createClass(ApiClient, [{
    key: 'op_start',
    value: function () {
      var _ref2 = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee(opStartReq) {
        return regeneratorRuntime.wrap(function _callee$(_context) {
          while (1) {
            switch (_context.prev = _context.next) {
              case 0:
                return _context.abrupt('return', fetch(this.baseUrl + '/ops/starts', {
                  method: 'POST',
                  body: JSON.stringify(opStartReq)
                }).then(function (response) {
                  return response.text();
                }));

              case 1:
              case 'end':
                return _context.stop();
            }
          }
        }, _callee, this);
      }));

      function op_start(_x) {
        return _ref2.apply(this, arguments);
      }

      return op_start;
    }()

    /**
     * Kills an op
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L139
     * @param opKillReq
     * @return {Promise.<null>}
     */

  }, {
    key: 'op_kill',
    value: function () {
      var _ref3 = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee2(opKillReq) {
        return regeneratorRuntime.wrap(function _callee2$(_context2) {
          while (1) {
            switch (_context2.prev = _context2.next) {
              case 0:
                return _context2.abrupt('return', fetch(this.baseUrl + '/ops/kills', {
                  method: 'POST',
                  body: JSON.stringify(opKillReq)
                }).then(function () {
                  return null;
                }));

              case 1:
              case 'end':
                return _context2.stop();
            }
          }
        }, _callee2, this);
      }));

      function op_kill(_x2) {
        return _ref3.apply(this, arguments);
      }

      return op_kill;
    }()

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
    value: function () {
      var _ref5 = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee3(_ref4) {
        var pkgRef = _ref4.pkgRef,
            contentPath = _ref4.contentPath;
        return regeneratorRuntime.wrap(function _callee3$(_context3) {
          while (1) {
            switch (_context3.prev = _context3.next) {
              case 0:
                return _context3.abrupt('return', fetch(this.baseUrl + '/pkgs/' + encodeURIComponent(pkgRef) + '/contents/' + encodeURIComponent(contentPath)));

              case 1:
              case 'end':
                return _context3.stop();
            }
          }
        }, _callee3, this);
      }));

      function pkg_content_get(_x3) {
        return _ref5.apply(this, arguments);
      }

      return pkg_content_get;
    }()

    /**
     * Lists pkg contents
     *
     * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L178
     * @param pkgRef
     * @return {Promise.<Object[]>}}
     */

  }, {
    key: 'pkg_content_list',
    value: function () {
      var _ref7 = _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee4(_ref6) {
        var pkgRef = _ref6.pkgRef;
        return regeneratorRuntime.wrap(function _callee4$(_context4) {
          while (1) {
            switch (_context4.prev = _context4.next) {
              case 0:
                return _context4.abrupt('return', fetch(this.baseUrl + '/pkgs/' + encodeURIComponent(pkgRef) + '/contents').then(function (response) {
                  return response.json();
                }));

              case 1:
              case 'end':
                return _context4.stop();
            }
          }
        }, _callee4, this);
      }));

      function pkg_content_list(_x4) {
        return _ref7.apply(this, arguments);
      }

      return pkg_content_list;
    }()
  }]);

  return ApiClient;
}();

module.exports = ApiClient;