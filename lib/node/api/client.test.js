const ObjectUnderTest = require('./client')
const Duplex = require('stream').Duplex

const providedBaseUrl = 'http://dummyBaseUrl'
const objectUnderTest = new ObjectUnderTest({baseUrl: providedBaseUrl})

const nock = require('nock')

afterEach(() => {
  nock.cleanAll()
})

describe('livenessGet', () => {
  it('makes expected http request', async () => {
    /* arrange */
    const scope = nock(providedBaseUrl)
      .log(console.log)
      .get(
        `/liveness`)
      .reply('200')

    /* act */
    await objectUnderTest.livenessGet()

    /* assert */
    scope.isDone()
  })
  describe('response is 300', async () => {
    it('returns expected result', async () => {
      /* arrange */
      const expectedErrMsg = 'dummyErrorMsg'

      nock(providedBaseUrl)
        .log(console.log)
        .get(
          `/liveness`)
        .reply('300', expectedErrMsg)

      /* act */
      let actualErr
      try {
        await objectUnderTest.livenessGet()
      } catch (err) {
        actualErr = err
      }

      /* assert */
      expect(actualErr.message).toEqual(expectedErrMsg)
    })
  })
  describe('response is 200', async () => {
    it('returns expected result', async () => {
      /* arrange */
      nock(providedBaseUrl)
        .log(console.log)
        .get(
          `/liveness`)
        .reply('200')

      /* act */
      const actualResponse = await objectUnderTest.livenessGet()

      /* assert */
      expect(actualResponse).toEqual('')
    })
  })
})

describe('opKill', () => {
  it('makes expected http request', async () => {
    /* arrange */
    const providedOpKillReq = {
      'opId': 'dummyOpId'
    }

    const scope = nock(providedBaseUrl)
      .log(console.log)
      .post('/ops/kills', JSON.stringify(providedOpKillReq))
      .reply('200', {})

    /* act */
    await objectUnderTest.opKill(providedOpKillReq)

    /* assert */
    scope.isDone()
  })
  describe('response is 300', async () => {
    it('returns expected result', async () => {
      /* arrange */
      const expectedErrMsg = 'dummyErrorMsg'

      nock(providedBaseUrl)
        .log(console.log)
        .post(
          `/ops/kills`)
        .reply('300', expectedErrMsg)

      /* act */
      let actualErr
      try {
        await objectUnderTest.opKill({})
      } catch (err) {
        actualErr = err
      }

      /* assert */
      expect(actualErr.message).toEqual(expectedErrMsg)
    })
  })
  describe('response is 202', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedOpKillReq = {
        'dummyName': 'dummyValue'
      }

      nock(providedBaseUrl)
        .log(console.log)
        .post('/ops/kills', JSON.stringify(providedOpKillReq))
        .reply('204')

      /* act */
      const actualResult = await objectUnderTest.opKill(providedOpKillReq)

      /* assert */
      expect(actualResult).toBeNull()
    })
  })
})

describe('opStart', () => {
  it('makes expected http request', async () => {
    /* arrange */
    const providedOpStartReq = {
      'dummyName': 'dummyValue'
    }

    const scope = nock(providedBaseUrl)
      .log(console.log)
      .post('/ops/starts', JSON.stringify(providedOpStartReq))
      .reply('200')

    /* act */
    await objectUnderTest.opStart(providedOpStartReq)

    /* assert */
    scope.isDone()
  })
  describe('response is 300', async () => {
    it('returns expected result', async () => {
      /* arrange */
      const expectedErrMsg = 'dummyErrorMsg'

      nock(providedBaseUrl)
        .log(console.log)
        .post(
          `/ops/starts`)
        .reply('300', expectedErrMsg)

      /* act */
      let actualErr
      try {
        await objectUnderTest.opStart({})
      } catch (err) {
        actualErr = err
      }

      /* assert */
      expect(actualErr.message).toEqual(expectedErrMsg)
    })
  })
  describe('response is 201', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedOpStartReq = {
        'dummyName': 'dummyValue'
      }

      const expectedOpId = 'dummyOpId'

      nock(providedBaseUrl)
        .log(console.log)
        .post('/ops/starts', JSON.stringify(providedOpStartReq))
        .reply('200', expectedOpId)

      /* act */
      const actualOpId = await objectUnderTest.opStart(providedOpStartReq)

      /* assert */
      expect(actualOpId).toEqual(expectedOpId)
    })
  })
})
describe('dataGet', () => {
  it('makes expected http request', async () => {
    /* arrange */
    const providedDataRef = '//dummyDataRef'

    const scope = nock(providedBaseUrl)
      .log(console.log)
      .get(
        `/data/${encodeURIComponent(providedDataRef)}`)
      .reply('200')

    /* act */
    await objectUnderTest.dataGet({
      dataRef: providedDataRef
    })

    /* assert */
    scope.isDone()
  })
  describe('response is 300', async () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedDataRef = '//dummyDataRef'
      const expectedErrMsg = 'dummyErrorMsg'

      nock(providedBaseUrl)
        .log(console.log)
        .get(
          `/data/${encodeURIComponent(providedDataRef)}`)
        .reply('300', expectedErrMsg)

      /* act */
      let actualErr
      try {
        await objectUnderTest.dataGet({
          dataRef: providedDataRef
        })
      } catch (err) {
        actualErr = err
      }

      /* assert */
      expect(actualErr.message).toEqual(expectedErrMsg)
    })
  })
  describe('response is 200', async () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedDataRef = '//dummyDataRef'

      const expectedContent = 'dummyContent'

      nock(providedBaseUrl)
        .log(console.log)
        .get(
          `/data/${encodeURIComponent(providedDataRef)}`, () => {
            let stream = new Duplex()
            stream.push(Buffer.from(expectedContent))
            stream.push(null)
            return stream
          })
        .reply('200', expectedContent)

      /* act */
      const actualContent = await objectUnderTest.dataGet({
        dataRef: providedDataRef
      }).then(response => response.text())

      /* assert */
      expect(actualContent).toEqual(expectedContent)
    })
  })
})
