const ObjectUnderTest = require('./apiClient');
const Duplex = require('stream').Duplex;

const providedBaseUrl = 'http://dummyBaseUrl';
const objectUnderTest = new ObjectUnderTest({baseUrl: providedBaseUrl});

const nock = require('nock');

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
    scope.isDone();
  });
  describe('response is 202', () => {
    test('returns expected result', async () => {
      /* arrange */
      const providedOpKillReq = {
        'dummyName': 'dummyValue'
      };

      nock(providedBaseUrl)
        .log(console.log)
        .post('/ops/kills', JSON.stringify(providedOpKillReq))
        .reply('204');

      /* act */
      const actualResult = await objectUnderTest.op_kill(providedOpKillReq);

      /* assert */
      expect(actualResult).toBeNull();
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
      .reply('200');

    /* act */
    await objectUnderTest.op_start(providedOpStartReq);

    /* assert */
    scope.isDone();
  });
  describe('response is 201', () => {
    test('returns expected result', async () => {
      /* arrange */
      const providedOpStartReq = {
        'dummyName': 'dummyValue'
      };

      const expectedOpId = 'dummyOpId';

      nock(providedBaseUrl)
        .log(console.log)
        .post('/ops/starts', JSON.stringify(providedOpStartReq))
        .reply('200', expectedOpId);

      /* act */
      const actualOpId = await objectUnderTest.op_start(providedOpStartReq);

      /* assert */
      expect(actualOpId).toEqual(expectedOpId);
    });
  })
});

describe('pkg_content_get', () => {
  test('makes expected http request', async () => {
    /* arrange */
    const providedPkgRef = '//dummyPkgRef';
    const providedContentPath = '/dummyContentPath';

    const scope = nock(providedBaseUrl)
      .log(console.log)
      .get(
        `/pkgs/${encodeURIComponent(providedPkgRef)}/contents/${encodeURIComponent(providedContentPath)}`)
      .reply('200');

    /* act */
    await objectUnderTest.pkg_content_get({
      pkgRef: providedPkgRef,
      contentPath: providedContentPath
    });

    /* assert */
    scope.isDone();
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
          `/pkgs/${encodeURIComponent(providedPkgRef)}/contents/${encodeURIComponent(providedContentPath)}`, () => {
            let stream = new Duplex();
            stream.push(Buffer.from(expectedContent));
            stream.push(null);
            return stream;
          })
        .reply('200', expectedContent);

      /* act */
      const actualContent = await objectUnderTest.pkg_content_get({
        pkgRef: providedPkgRef,
        contentPath: providedContentPath
      }).then(response => response.text());

      /* assert */
      expect(actualContent).toEqual(expectedContent)
    });
  })
});

describe('pkg_content_list', () => {
  test('makes expected http request', async () => {
    /* arrange */
    const providedPkgRef = '//dummyPkgRef';

    const scope = nock(providedBaseUrl)
      .log(console.log)
      .get(
        `/pkgs/${encodeURIComponent(providedPkgRef)}/contents`)
      .reply('200',[]);

    /* act */
    await objectUnderTest.pkg_content_list({pkgRef: providedPkgRef});

    /* assert */
    scope.isDone();
  });
  describe('response is 200', async () => {
    test('returns expected result', async () => {
      /* arrange */
      const providedPkgRef = '//dummyPkgRef';

      const expectedContentList = [{path: '/dummy/pkg/path'}];

      nock(providedBaseUrl)
        .log(console.log)
        .get(
          `/pkgs/${encodeURIComponent(providedPkgRef)}/contents`)
        .reply('200', expectedContentList);

      /* act */
      const actualContentList = await objectUnderTest.pkg_content_list({pkgRef: providedPkgRef});

      /* assert */
      expect(actualContentList).toEqual(expectedContentList)
    });
  })
});
