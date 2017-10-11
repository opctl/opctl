'use strict';

var ObjectUnderTest = require('./apiClient');
var Duplex = require('stream').Duplex;

var providedBaseUrl = 'http://dummyBaseUrl';
var objectUnderTest = new ObjectUnderTest({ baseUrl: providedBaseUrl });

var nock = require('nock');

afterEach(function () {
  nock.cleanAll();
});

describe('op_kill', function () {
  test('makes expected http request', async function () {
    /* arrange */
    var providedOpKillReq = {
      'opId': 'dummyOpId'
    };

    var scope = nock(providedBaseUrl).log(console.log).post('/ops/kills', JSON.stringify(providedOpKillReq)).reply('200', {});

    /* act */
    await objectUnderTest.op_kill(providedOpKillReq);

    /* assert */
    scope.isDone();
  });
  describe('response is 202', function () {
    test('returns expected result', async function () {
      /* arrange */
      var providedOpKillReq = {
        'dummyName': 'dummyValue'
      };

      nock(providedBaseUrl).log(console.log).post('/ops/kills', JSON.stringify(providedOpKillReq)).reply('204');

      /* act */
      var actualResult = await objectUnderTest.op_kill(providedOpKillReq);

      /* assert */
      expect(actualResult).toBeNull();
    });
  });
});

describe('op_start', function () {
  test('makes expected http request', async function () {
    /* arrange */
    var providedOpStartReq = {
      'dummyName': 'dummyValue'
    };

    var scope = nock(providedBaseUrl).log(console.log).post('/ops/starts', JSON.stringify(providedOpStartReq)).reply('200');

    /* act */
    await objectUnderTest.op_start(providedOpStartReq);

    /* assert */
    scope.isDone();
  });
  describe('response is 201', function () {
    test('returns expected result', async function () {
      /* arrange */
      var providedOpStartReq = {
        'dummyName': 'dummyValue'
      };

      var expectedOpId = 'dummyOpId';

      nock(providedBaseUrl).log(console.log).post('/ops/starts', JSON.stringify(providedOpStartReq)).reply('200', expectedOpId);

      /* act */
      var actualOpId = await objectUnderTest.op_start(providedOpStartReq);

      /* assert */
      expect(actualOpId).toEqual(expectedOpId);
    });
  });
});

describe('pkg_content_get', function () {
  test('makes expected http request', async function () {
    /* arrange */
    var providedPkgRef = '//dummyPkgRef';
    var providedContentPath = '/dummyContentPath';

    var scope = nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents/' + encodeURIComponent(providedContentPath)).reply('200');

    /* act */
    await objectUnderTest.pkg_content_get({
      pkgRef: providedPkgRef,
      contentPath: providedContentPath
    });

    /* assert */
    scope.isDone();
  });
  describe('response is 200', async function () {
    test('returns expected result', async function () {
      /* arrange */
      var providedPkgRef = '//dummyPkgRef';
      var providedContentPath = '/dummyContentPath';

      var expectedContent = 'dummyContent';

      nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents/' + encodeURIComponent(providedContentPath), function () {
        var stream = new Duplex();
        stream.push(Buffer.from(expectedContent));
        stream.push(null);
        return stream;
      }).reply('200', expectedContent);

      /* act */
      var actualContent = await objectUnderTest.pkg_content_get({
        pkgRef: providedPkgRef,
        contentPath: providedContentPath
      }).then(function (response) {
        return response.text();
      });

      /* assert */
      expect(actualContent).toEqual(expectedContent);
    });
  });
});

describe('pkg_content_list', function () {
  test('makes expected http request', async function () {
    /* arrange */
    var providedPkgRef = '//dummyPkgRef';

    var scope = nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents').reply('200', []);

    /* act */
    await objectUnderTest.pkg_content_list({ pkgRef: providedPkgRef });

    /* assert */
    scope.isDone();
  });
  describe('response is 200', async function () {
    test('returns expected result', async function () {
      /* arrange */
      var providedPkgRef = '//dummyPkgRef';

      var expectedContentList = [{ path: '/dummy/pkg/path' }];

      nock(providedBaseUrl).log(console.log).get('/pkgs/' + encodeURIComponent(providedPkgRef) + '/contents').reply('200', expectedContentList);

      /* act */
      var actualContentList = await objectUnderTest.pkg_content_list({ pkgRef: providedPkgRef });

      /* assert */
      expect(actualContentList).toEqual(expectedContentList);
    });
  });
});