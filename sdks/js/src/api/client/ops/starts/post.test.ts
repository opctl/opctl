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
    expect(scope.isDone()).toBe(true)
  })
  it.each([undefined, false])('omits noCache when %s', async (noCache) => {
    const args = {}
    const op = { ref: 'ref' }
    const scope = nock(providedNetRef)
      .post('/ops/starts', {
        args,
        op
      })
      .reply(200, 'opId')

    expect(await objectUnderTest(providedNetRef, args, op, noCache)).toEqual('opId')
    expect(scope.isDone()).toBe(true)
  })
  it('sends noCache when true', async () => {
    const scope = nock(providedNetRef)
      .post('/ops/starts', {
        args: {},
        op: { ref: 'ref' },
        noCache: true
      })
      .reply(200, 'opId')

    expect(await objectUnderTest(providedNetRef, {}, { ref: 'ref' }, true)).toEqual('opId')
    expect(scope.isDone()).toBe(true)
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
