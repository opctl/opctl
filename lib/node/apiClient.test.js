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

describe('liveness_get', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee() {
    var scope;
    return regeneratorRuntime.wrap(function _callee$(_context) {
      while (1) {
        switch (_context.prev = _context.next) {
          case 0:
            /* arrange */
            scope = nock(providedBaseUrl).log(console.log).get('/liveness').reply('200');

            /* act */

            _context.next = 3;
            return objectUnderTest.liveness_get();

          case 3:

            /* assert */
            scope.isDone();

          case 4:
          case 'end':
            return _context.stop();
        }
      }
    }, _callee, undefined);
  })));
  describe('response is 300', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee3() {
    return regeneratorRuntime.wrap(function _callee3$(_context3) {
      while (1) {
        switch (_context3.prev = _context3.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee2() {
              var expectedErrMsg, actualErr;
              return regeneratorRuntime.wrap(function _callee2$(_context2) {
                while (1) {
                  switch (_context2.prev = _context2.next) {
                    case 0:
                      /* arrange */
                      expectedErrMsg = 'dummyErrorMsg';


                      nock(providedBaseUrl).log(console.log).get('/liveness').reply('300', expectedErrMsg);

                      /* act */
                      actualErr = void 0;
                      _context2.prev = 3;
                      _context2.next = 6;
                      return objectUnderTest.liveness_get();

                    case 6:
                      _context2.next = 11;
                      break;

                    case 8:
                      _context2.prev = 8;
                      _context2.t0 = _context2['catch'](3);

                      actualErr = _context2.t0;

                    case 11:

                      /* assert */
                      expect(actualErr.message).toEqual(expectedErrMsg);

                    case 12:
                    case 'end':
                      return _context2.stop();
                  }
                }
              }, _callee2, undefined, [[3, 8]]);
            })));

          case 1:
          case 'end':
            return _context3.stop();
        }
      }
    }, _callee3, undefined);
  })));
  describe('response is 200', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee5() {
    return regeneratorRuntime.wrap(function _callee5$(_context5) {
      while (1) {
        switch (_context5.prev = _context5.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee4() {
              var actualResponse;
              return regeneratorRuntime.wrap(function _callee4$(_context4) {
                while (1) {
                  switch (_context4.prev = _context4.next) {
                    case 0:
                      /* arrange */
                      nock(providedBaseUrl).log(console.log).get('/liveness').reply('200');

                      /* act */
                      _context4.next = 3;
                      return objectUnderTest.liveness_get();

                    case 3:
                      actualResponse = _context4.sent;


                      /* assert */
                      expect(actualResponse).toEqual('');

                    case 5:
                    case 'end':
                      return _context4.stop();
                  }
                }
              }, _callee4, undefined);
            })));

          case 1:
          case 'end':
            return _context5.stop();
        }
      }
    }, _callee5, undefined);
  })));
});

describe('op_kill', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee6() {
    var providedOpKillReq, scope;
    return regeneratorRuntime.wrap(function _callee6$(_context6) {
      while (1) {
        switch (_context6.prev = _context6.next) {
          case 0:
            /* arrange */
            providedOpKillReq = {
              'opId': 'dummyOpId'
            };
            scope = nock(providedBaseUrl).log(console.log).post('/ops/kills', JSON.stringify(providedOpKillReq)).reply('200', {});

            /* act */

            _context6.next = 4;
            return objectUnderTest.op_kill(providedOpKillReq);

          case 4:

            /* assert */
            scope.isDone();

          case 5:
          case 'end':
            return _context6.stop();
        }
      }
    }, _callee6, undefined);
  })));
  describe('response is 300', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee8() {
    return regeneratorRuntime.wrap(function _callee8$(_context8) {
      while (1) {
        switch (_context8.prev = _context8.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee7() {
              var expectedErrMsg, actualErr;
              return regeneratorRuntime.wrap(function _callee7$(_context7) {
                while (1) {
                  switch (_context7.prev = _context7.next) {
                    case 0:
                      /* arrange */
                      expectedErrMsg = 'dummyErrorMsg';


                      nock(providedBaseUrl).log(console.log).post('/ops/kills').reply('300', expectedErrMsg);

                      /* act */
                      actualErr = void 0;
                      _context7.prev = 3;
                      _context7.next = 6;
                      return objectUnderTest.op_kill({});

                    case 6:
                      _context7.next = 11;
                      break;

                    case 8:
                      _context7.prev = 8;
                      _context7.t0 = _context7['catch'](3);

                      actualErr = _context7.t0;

                    case 11:

                      /* assert */
                      expect(actualErr.message).toEqual(expectedErrMsg);

                    case 12:
                    case 'end':
                      return _context7.stop();
                  }
                }
              }, _callee7, undefined, [[3, 8]]);
            })));

          case 1:
          case 'end':
            return _context8.stop();
        }
      }
    }, _callee8, undefined);
  })));
  describe('response is 202', function () {
    test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee9() {
      var providedOpKillReq, actualResult;
      return regeneratorRuntime.wrap(function _callee9$(_context9) {
        while (1) {
          switch (_context9.prev = _context9.next) {
            case 0:
              /* arrange */
              providedOpKillReq = {
                'dummyName': 'dummyValue'
              };


              nock(providedBaseUrl).log(console.log).post('/ops/kills', JSON.stringify(providedOpKillReq)).reply('204');

              /* act */
              _context9.next = 4;
              return objectUnderTest.op_kill(providedOpKillReq);

            case 4:
              actualResult = _context9.sent;


              /* assert */
              expect(actualResult).toBeNull();

            case 6:
            case 'end':
              return _context9.stop();
          }
        }
      }, _callee9, undefined);
    })));
  });
});

describe('op_start', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee10() {
    var providedOpStartReq, scope;
    return regeneratorRuntime.wrap(function _callee10$(_context10) {
      while (1) {
        switch (_context10.prev = _context10.next) {
          case 0:
            /* arrange */
            providedOpStartReq = {
              'dummyName': 'dummyValue'
            };
            scope = nock(providedBaseUrl).log(console.log).post('/ops/starts', JSON.stringify(providedOpStartReq)).reply('200');

            /* act */

            _context10.next = 4;
            return objectUnderTest.op_start(providedOpStartReq);

          case 4:

            /* assert */
            scope.isDone();

          case 5:
          case 'end':
            return _context10.stop();
        }
      }
    }, _callee10, undefined);
  })));
  describe('response is 300', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee12() {
    return regeneratorRuntime.wrap(function _callee12$(_context12) {
      while (1) {
        switch (_context12.prev = _context12.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee11() {
              var expectedErrMsg, actualErr;
              return regeneratorRuntime.wrap(function _callee11$(_context11) {
                while (1) {
                  switch (_context11.prev = _context11.next) {
                    case 0:
                      /* arrange */
                      expectedErrMsg = 'dummyErrorMsg';


                      nock(providedBaseUrl).log(console.log).post('/ops/starts').reply('300', expectedErrMsg);

                      /* act */
                      actualErr = void 0;
                      _context11.prev = 3;
                      _context11.next = 6;
                      return objectUnderTest.op_start({});

                    case 6:
                      _context11.next = 11;
                      break;

                    case 8:
                      _context11.prev = 8;
                      _context11.t0 = _context11['catch'](3);

                      actualErr = _context11.t0;

                    case 11:

                      /* assert */
                      expect(actualErr.message).toEqual(expectedErrMsg);

                    case 12:
                    case 'end':
                      return _context11.stop();
                  }
                }
              }, _callee11, undefined, [[3, 8]]);
            })));

          case 1:
          case 'end':
            return _context12.stop();
        }
      }
    }, _callee12, undefined);
  })));
  describe('response is 201', function () {
    test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee13() {
      var providedOpStartReq, expectedOpId, actualOpId;
      return regeneratorRuntime.wrap(function _callee13$(_context13) {
        while (1) {
          switch (_context13.prev = _context13.next) {
            case 0:
              /* arrange */
              providedOpStartReq = {
                'dummyName': 'dummyValue'
              };
              expectedOpId = 'dummyOpId';


              nock(providedBaseUrl).log(console.log).post('/ops/starts', JSON.stringify(providedOpStartReq)).reply('200', expectedOpId);

              /* act */
              _context13.next = 5;
              return objectUnderTest.op_start(providedOpStartReq);

            case 5:
              actualOpId = _context13.sent;


              /* assert */
              expect(actualOpId).toEqual(expectedOpId);

            case 7:
            case 'end':
              return _context13.stop();
          }
        }
      }, _callee13, undefined);
    })));
  });
});

describe('pkg_content_get', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee14() {
    var providedPkgRef, providedContentPath, scope;
    return regeneratorRuntime.wrap(function _callee14$(_context14) {
      while (1) {
        switch (_context14.prev = _context14.next) {
          case 0:
            /* arrange */
            providedPkgRef = '//dummyPkgRef';
            providedContentPath = '/dummyContentPath';
            scope = nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents/' + encodeURIComponent(providedContentPath)).reply('200');

            /* act */

            _context14.next = 5;
            return objectUnderTest.pkg_content_get({
              pkgRef: providedPkgRef,
              contentPath: providedContentPath
            });

          case 5:

            /* assert */
            scope.isDone();

          case 6:
          case 'end':
            return _context14.stop();
        }
      }
    }, _callee14, undefined);
  })));
  describe('response is 300', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee16() {
    return regeneratorRuntime.wrap(function _callee16$(_context16) {
      while (1) {
        switch (_context16.prev = _context16.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee15() {
              var providedPkgRef, providedContentPath, expectedErrMsg, actualErr;
              return regeneratorRuntime.wrap(function _callee15$(_context15) {
                while (1) {
                  switch (_context15.prev = _context15.next) {
                    case 0:
                      /* arrange */
                      providedPkgRef = '//dummyPkgRef';
                      providedContentPath = '/dummyContentPath';
                      expectedErrMsg = 'dummyErrorMsg';


                      nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents/' + encodeURIComponent(providedContentPath)).reply('300', expectedErrMsg);

                      /* act */
                      actualErr = void 0;
                      _context15.prev = 5;
                      _context15.next = 8;
                      return objectUnderTest.pkg_content_get({
                        pkgRef: providedPkgRef,
                        contentPath: providedContentPath
                      });

                    case 8:
                      _context15.next = 13;
                      break;

                    case 10:
                      _context15.prev = 10;
                      _context15.t0 = _context15['catch'](5);

                      actualErr = _context15.t0;

                    case 13:

                      /* assert */
                      expect(actualErr.message).toEqual(expectedErrMsg);

                    case 14:
                    case 'end':
                      return _context15.stop();
                  }
                }
              }, _callee15, undefined, [[5, 10]]);
            })));

          case 1:
          case 'end':
            return _context16.stop();
        }
      }
    }, _callee16, undefined);
  })));
  describe('response is 200', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee18() {
    return regeneratorRuntime.wrap(function _callee18$(_context18) {
      while (1) {
        switch (_context18.prev = _context18.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee17() {
              var providedPkgRef, providedContentPath, expectedContent, actualContent;
              return regeneratorRuntime.wrap(function _callee17$(_context17) {
                while (1) {
                  switch (_context17.prev = _context17.next) {
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
                      _context17.next = 6;
                      return objectUnderTest.pkg_content_get({
                        pkgRef: providedPkgRef,
                        contentPath: providedContentPath
                      }).then(function (response) {
                        return response.text();
                      });

                    case 6:
                      actualContent = _context17.sent;


                      /* assert */
                      expect(actualContent).toEqual(expectedContent);

                    case 8:
                    case 'end':
                      return _context17.stop();
                  }
                }
              }, _callee17, undefined);
            })));

          case 1:
          case 'end':
            return _context18.stop();
        }
      }
    }, _callee18, undefined);
  })));
});

describe('pkg_content_list', function () {
  test('makes expected http request', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee19() {
    var providedPkgRef, scope;
    return regeneratorRuntime.wrap(function _callee19$(_context19) {
      while (1) {
        switch (_context19.prev = _context19.next) {
          case 0:
            /* arrange */
            providedPkgRef = '//dummyPkgRef';
            scope = nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents').reply('200', []);

            /* act */

            _context19.next = 4;
            return objectUnderTest.pkg_content_list({ pkgRef: providedPkgRef });

          case 4:

            /* assert */
            scope.isDone();

          case 5:
          case 'end':
            return _context19.stop();
        }
      }
    }, _callee19, undefined);
  })));
  describe('response is 300', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee21() {
    return regeneratorRuntime.wrap(function _callee21$(_context21) {
      while (1) {
        switch (_context21.prev = _context21.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee20() {
              var providedPkgRef, expectedErrMsg, actualErr;
              return regeneratorRuntime.wrap(function _callee20$(_context20) {
                while (1) {
                  switch (_context20.prev = _context20.next) {
                    case 0:
                      /* arrange */
                      providedPkgRef = '//dummyPkgRef';
                      expectedErrMsg = 'dummyErrorMsg';


                      nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents').reply('300', expectedErrMsg);

                      /* act */
                      actualErr = void 0;
                      _context20.prev = 4;
                      _context20.next = 7;
                      return objectUnderTest.pkg_content_list({
                        pkgRef: providedPkgRef
                      });

                    case 7:
                      _context20.next = 12;
                      break;

                    case 9:
                      _context20.prev = 9;
                      _context20.t0 = _context20['catch'](4);

                      actualErr = _context20.t0;

                    case 12:

                      /* assert */
                      expect(actualErr.message).toEqual(expectedErrMsg);

                    case 13:
                    case 'end':
                      return _context20.stop();
                  }
                }
              }, _callee20, undefined, [[4, 9]]);
            })));

          case 1:
          case 'end':
            return _context21.stop();
        }
      }
    }, _callee21, undefined);
  })));
  describe('response is 200', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee23() {
    return regeneratorRuntime.wrap(function _callee23$(_context23) {
      while (1) {
        switch (_context23.prev = _context23.next) {
          case 0:
            test('returns expected result', _asyncToGenerator( /*#__PURE__*/regeneratorRuntime.mark(function _callee22() {
              var providedPkgRef, expectedContentList, actualContentList;
              return regeneratorRuntime.wrap(function _callee22$(_context22) {
                while (1) {
                  switch (_context22.prev = _context22.next) {
                    case 0:
                      /* arrange */
                      providedPkgRef = '//dummyPkgRef';
                      expectedContentList = [{ path: '/dummy/pkg/path' }];


                      nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents').reply('200', expectedContentList);

                      /* act */
                      _context22.next = 5;
                      return objectUnderTest.pkg_content_list({ pkgRef: providedPkgRef });

                    case 5:
                      actualContentList = _context22.sent;


                      /* assert */
                      expect(actualContentList).toEqual(expectedContentList);

                    case 7:
                    case 'end':
                      return _context22.stop();
                  }
                }
              }, _callee22, undefined);
            })));

          case 1:
          case 'end':
            return _context23.stop();
        }
      }
    }, _callee23, undefined);
  })));
});