import objectUnderTest from './post'

const providedNetRef = 'http://dummyNetRef'

const nock = require('nock')

afterEach(() => {
  nock.cleanAll()
})

describe('opKill', () => {
  it('makes expected http request', async () => {
    /* arrange */
    const providedOpId = 'providedOpId'

    const scope = nock(providedNetRef)
      
      .post('/ops/kills', JSON.stringify({
        opId: providedOpId
      }))
      .reply(200, {})

    /* act */
    await objectUnderTest(
      providedNetRef,
      providedOpId
    )

    /* assert */
    scope.isDone()
  })
  describe('response is 300', () => {
    it('returns expected result', async () => {
      /* arrange */
      const expectedErrMsg = 'dummyErrorMsg'

      nock(providedNetRef)
        
        .post(
          `/ops/kills`)
        .reply(300, expectedErrMsg)

      /* act */
      let actualErr: null| Error = null
      try {
        await objectUnderTest(
          providedNetRef,
          'opId'
        )
      } catch (err) {
        actualErr = err as Error
      }

      /* assert */
      expect(actualErr?.message).toEqual(expectedErrMsg)
    })
  })
  describe('response is 202', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedOpKillReq = {
        'dummyName': 'dummyValue'
      }

      nock(providedNetRef)
        
        .post('/ops/kills')
        .reply(204)

      /* act */
      const actualResult = await objectUnderTest(
        providedNetRef,
        'opId'
      )

      /* assert */
      expect(actualResult).toBeNull()
    })
  })
})
