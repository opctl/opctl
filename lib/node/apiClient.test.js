const ObjectUnderTest = require('./apiClient');

const providedBaseUrl = 'http://dummyBaseUrl';
const objectUnderTest = new ObjectUnderTest({baseUrl: providedBaseUrl});

const nock = require('nock');
const axios = require('axios');

axios.defaults.adapter = require('axios/lib/adapters/http');

afterEach(() => {
  nock.cleanAll();
});

describe('op_kill', () => {
  test('makes expected http request', async () => {
    /* arrange */
    const providedOpKillReq = {
      'opId': 'dummyOpId'
    };

    const scope = nock(providedBaseUrl)
      .log(console.log)
      .post('/ops/kills', JSON.stringify(providedOpKillReq))
      .reply('200', {});

    /* act */
    await objectUnderTest.op_kill(providedOpKillReq);

    /* assert */
    scope.done();
  });
  describe('response is 202', () => {
    test('returns expected result', async () => {
      /* arrange */
      const providedOpKillReq = {
        'dummyName': 'dummyValue'
      };

      const scope = nock(providedBaseUrl)
        .log(console.log)
        .post('/ops/kills', JSON.stringify(providedOpKillReq))
        .reply('204');

      /* act */
      await objectUnderTest.op_kill(providedOpKillReq);

      /* assert */
      scope.done();
    });
  })
});

describe('op_start', () => {
  test('makes expected http request', async () => {
    /* arrange */
    const providedOpStartReq = {
      'dummyName': 'dummyValue'
    };

    const scope = nock(providedBaseUrl)
      .log(console.log)
      .post('/ops/starts', JSON.stringify(providedOpStartReq))
      .reply('200', {});

    /* act */
    await objectUnderTest.op_start(providedOpStartReq);

    /* assert */
    scope.done();
  });
  describe('response is 201', () => {
    test('returns expected result', async () => {
      /* arrange */
      const providedOpStartReq = {
        'dummyName': 'dummyValue'
      };

      const scope = nock(providedBaseUrl)
        .log(console.log)
        .post('/ops/starts', JSON.stringify(providedOpStartReq))
        .reply('200', {});

      /* act */
      await objectUnderTest.op_start(providedOpStartReq);

      /* assert */
      scope.done();
    });
  })
});

describe('pkg_content_get', () => {
  test('makes expected http request', async () => {
    /* arrange */
    const providedPkgRef = '//dummyPkgRef';
    const providedContentPath = '/dummyContentPath';

    const expectedContent = 'dummyContent';

    nock(providedBaseUrl)
      .log(console.log)
      .get(
        `/pkgs/${encodeURIComponent(providedPkgRef)}/contents/${encodeURIComponent(providedContentPath)}`)
      .reply('200', expectedContent);

    /* act */
    const actualOpId = await objectUnderTest.pkg_content_get({pkgRef: providedPkgRef, contentPath: providedContentPath});

    /* assert */
    expect(actualOpId).toEqual(actualOpId)
  });
  describe('response is 200', async () => {
    test('returns expected result', async () => {
      /* arrange */
      const providedPkgRef = '//dummyPkgRef';
      const providedContentPath = '/dummyContentPath';

      const expectedContent = 'dummyContent';

      nock(providedBaseUrl)
        .log(console.log)
        .get(
          `/pkgs/${encodeURIComponent(providedPkgRef)}/contents/${encodeURIComponent(providedContentPath)}`)
        .reply('200', expectedContent);

      /* act */
      const actualContent = await objectUnderTest.pkg_content_get({pkgRef: providedPkgRef, contentPath: providedContentPath});

      /* assert */
      expect(actualContent).toEqual(expectedContent)
    });
  })
});
