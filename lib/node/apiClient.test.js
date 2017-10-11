'use strict';

function _asyncToGenerator(fn) { return function () { var gen = fn.apply(this, arguments); return new Promise(function (resolve, reject) { function step(key, arg) { try { var info = gen[key](arg); var value = info.value; } catch (error) { reject(error); return; } if (info.done) { resolve(value); } else { return Promise.resolve(value).then(function (value) { step("next", value); }, function (err) { step("throw", err); }); } } return step("next"); }); }; }

var ObjectUnderTest = require('./apiClient');
var Duplex = require('stream').Duplex;

var providedBaseUrl = 'http://dummyBaseUrl';
var objectUnderTest = new ObjectUnderTest({ baseUrl: providedBaseUrl });

var nock = require('nock');

afterEach(function () {
  nock.cleanAll();
});

describe('op_kill', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee() {
    var providedOpKillReq, scope;
    return regeneratorRuntime.wrap(function _callee$(_context) {
      while (1) {
        switch (_context.prev = _context.next) {
          case 0:
            /* arrange */
            providedOpKillReq = {
              'opId': 'dummyOpId'
            };
            scope = nock(providedBaseUrl).log(console.log).post('/ops/kills', JSON.stringify(providedOpKillReq)).reply('200', {});

            /* act */

            _context.next = 4;
            return objectUnderTest.op_kill(providedOpKillReq);

          case 4:

            /* assert */
            scope.isDone();

          case 5:
          case 'end':
            return _context.stop();
        }
      }
    }, _callee, undefined);
  })));
  describe('response is 202', function () {
    test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee2() {
      var providedOpKillReq, actualResult;
      return regeneratorRuntime.wrap(function _callee2$(_context2) {
        while (1) {
          switch (_context2.prev = _context2.next) {
            case 0:
              /* arrange */
              providedOpKillReq = {
                'dummyName': 'dummyValue'
              };


              nock(providedBaseUrl).log(console.log).post('/ops/kills', JSON.stringify(providedOpKillReq)).reply('204');

              /* act */
              _context2.next = 4;
              return objectUnderTest.op_kill(providedOpKillReq);

            case 4:
              actualResult = _context2.sent;


              /* assert */
              expect(actualResult).toBeNull();

            case 6:
            case 'end':
              return _context2.stop();
          }
        }
      }, _callee2, undefined);
    })));
  });
});

describe('op_start', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee3() {
    var providedOpStartReq, scope;
    return regeneratorRuntime.wrap(function _callee3$(_context3) {
      while (1) {
        switch (_context3.prev = _context3.next) {
          case 0:
            /* arrange */
            providedOpStartReq = {
              'dummyName': 'dummyValue'
            };
            scope = nock(providedBaseUrl).log(console.log).post('/ops/starts', JSON.stringify(providedOpStartReq)).reply('200');

            /* act */

            _context3.next = 4;
            return objectUnderTest.op_start(providedOpStartReq);

          case 4:

            /* assert */
            scope.isDone();

          case 5:
          case 'end':
            return _context3.stop();
        }
      }
    }, _callee3, undefined);
  })));
  describe('response is 201', function () {
    test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee4() {
      var providedOpStartReq, expectedOpId, actualOpId;
      return regeneratorRuntime.wrap(function _callee4$(_context4) {
        while (1) {
          switch (_context4.prev = _context4.next) {
            case 0:
              /* arrange */
              providedOpStartReq = {
                'dummyName': 'dummyValue'
              };
              expectedOpId = 'dummyOpId';


              nock(providedBaseUrl).log(console.log).post('/ops/starts', JSON.stringify(providedOpStartReq)).reply('200', expectedOpId);

              /* act */
              _context4.next = 5;
              return objectUnderTest.op_start(providedOpStartReq);

            case 5:
              actualOpId = _context4.sent;


              /* assert */
              expect(actualOpId).toEqual(expectedOpId);

            case 7:
            case 'end':
              return _context4.stop();
          }
        }
      }, _callee4, undefined);
    })));
  });
});

describe('pkg_content_get', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee5() {
    var providedPkgRef, providedContentPath, scope;
    return regeneratorRuntime.wrap(function _callee5$(_context5) {
      while (1) {
        switch (_context5.prev = _context5.next) {
          case 0:
            /* arrange */
            providedPkgRef = '//dummyPkgRef';
            providedContentPath = '/dummyContentPath';
            scope = nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents/' + encodeURIComponent(providedContentPath)).reply('200');

            /* act */

            _context5.next = 5;
            return objectUnderTest.pkg_content_get({
              pkgRef: providedPkgRef,
              contentPath: providedContentPath
            });

          case 5:

            /* assert */
            scope.isDone();

          case 6:
          case 'end':
            return _context5.stop();
        }
      }
    }, _callee5, undefined);
  })));
  describe('response is 200', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee7() {
    return regeneratorRuntime.wrap(function _callee7$(_context7) {
      while (1) {
        switch (_context7.prev = _context7.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee6() {
              var providedPkgRef, providedContentPath, expectedContent, actualContent;
              return regeneratorRuntime.wrap(function _callee6$(_context6) {
                while (1) {
                  switch (_context6.prev = _context6.next) {
                    case 0:
                      /* arrange */
                      providedPkgRef = '//dummyPkgRef';
                      providedContentPath = '/dummyContentPath';
                      expectedContent = 'dummyContent';


                      nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents/' + encodeURIComponent(providedContentPath), function () {
                        var stream = new Duplex();
                        stream.push(Buffer.from(expectedContent));
                        stream.push(null);
                        return stream;
                      }).reply('200', expectedContent);

                      /* act */
                      _context6.next = 6;
                      return objectUnderTest.pkg_content_get({
                        pkgRef: providedPkgRef,
                        contentPath: providedContentPath
                      }).then(function (response) {
                        return response.text();
                      });

                    case 6:
                      actualContent = _context6.sent;


                      /* assert */
                      expect(actualContent).toEqual(expectedContent);

                    case 8:
                    case 'end':
                      return _context6.stop();
                  }
                }
              }, _callee6, undefined);
            })));

          case 1:
          case 'end':
            return _context7.stop();
        }
      }
    }, _callee7, undefined);
  })));
});

describe('pkg_content_list', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee8() {
    var providedPkgRef, scope;
    return regeneratorRuntime.wrap(function _callee8$(_context8) {
      while (1) {
        switch (_context8.prev = _context8.next) {
          case 0:
            /* arrange */
            providedPkgRef = '//dummyPkgRef';
            scope = nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents').reply('200', []);

            /* act */

            _context8.next = 4;
            return objectUnderTest.pkg_content_list({ pkgRef: providedPkgRef });

          case 4:

            /* assert */
            scope.isDone();

          case 5:
          case 'end':
            return _context8.stop();
        }
      }
    }, _callee8, undefined);
  })));
  describe('response is 200', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee10() {
    return regeneratorRuntime.wrap(function _callee10$(_context10) {
      while (1) {
        switch (_context10.prev = _context10.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee9() {
              var providedPkgRef, expectedContentList, actualContentList;
              return regeneratorRuntime.wrap(function _callee9$(_context9) {
                while (1) {
                  switch (_context9.prev = _context9.next) {
                    case 0:
                      /* arrange */
                      providedPkgRef = '//dummyPkgRef';
                      expectedContentList = [{ path: '/dummy/pkg/path' }];


                      nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents').reply('200', expectedContentList);

                      /* act */
                      _context9.next = 5;
                      return objectUnderTest.pkg_content_list({ pkgRef: providedPkgRef });

                    case 5:
                      actualContentList = _context9.sent;


                      /* assert */
                      expect(actualContentList).toEqual(expectedContentList);

                    case 7:
                    case 'end':
                      return _context9.stop();
                  }
                }
              }, _callee9, undefined);
            })));

          case 1:
          case 'end':
            return _context10.stop();
        }
      }
    }, _callee10, undefined);
  })));
});