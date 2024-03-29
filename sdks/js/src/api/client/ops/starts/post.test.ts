import objectUnderTest from './post'

const nock = require('nock')
const providedNetRef = 'http://dummyNetRef'

afterEach(() => {
  nock.cleanAll()
})

describe('opStart', () => {
  it('makes expected http request', async () => {
    /* arrange */
    const providedArgs = {
      arg1Name: {
        string: 'string'
      }
    }
    const providedOp = {
      ref: 'ref'
    }

    const scope = nock(providedNetRef)
      
      .post(
        '/ops/starts',
        JSON.stringify({
          args: providedArgs,
          op: providedOp
        })
      )
      .reply(200)

    /* act */
    await objectUnderTest(
      providedNetRef,
      providedArgs,
      providedOp
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
          `/ops/starts`)
        .reply(300, expectedErrMsg)

      /* act */
      let actualErr: null | Error = null
      try {
        await objectUnderTest(
          providedNetRef,
          {
            ['']: {
              string: ''
            }
          },
          {
            ref: ''
          }
        )
      } catch (err) {
        actualErr = err as Error
      }

      /* assert */
      expect(actualErr?.message).toEqual(expectedErrMsg)
    })
  })
  describe('response is 201', () => {
    it('returns expected result', async () => {
      /* arrange */
      const expectedOpId = 'dummyOpId'

      nock(providedNetRef)
        
        .post('/ops/starts')
        .reply(200, expectedOpId)

      /* act */
      const actualOpId = await objectUnderTest(
        providedNetRef,
        {
          ['']: {
            string: ''
          }
        },
        {
          ref: ''
        }
      )

      /* assert */
      expect(actualOpId).toEqual(expectedOpId)
    })
  })
})
