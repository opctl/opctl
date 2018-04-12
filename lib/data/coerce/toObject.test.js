const nodeApiClient = require('../../node/api/client')

const objectUnderTest = require('./toObject')

describe('toObject', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('value is array', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { array: [] }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce array to object; incompatible types`)
    })
  })
  describe('value is boolean', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { boolean: true }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce boolean '${providedValue.boolean}' to object; incompatible types`)
    })
  })
  describe('value is dir', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { dir: 'dummyDir' }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce dir '${providedValue.dir}' to object; incompatible types`)
    })
  })
  describe('value is file', () => {
    it('should call nodeApiClient.dataGet w/ expected args', async () => {
      /* arrange */
      const providedValue = { file: 'dummyFile' }

      const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify({})) }
      nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

      /* act */
      await objectUnderTest(providedValue)

      /* assert */
      expect(nodeApiClient.dataGet)
        .toBeCalledWith({ dataRef: providedValue.file })
    })
    describe('nodeApiClient.dataGet errs', () => {
      it('returns expected result', async () => {
        /* arrange */
        const providedValue = { file: 'dummyFile' }

        const dataGetErr = new Error('dummyErr')
        nodeApiClient.dataGet = jest.fn().mockRejectedValue(dataGetErr)

        /* act/assert */
        await expect(objectUnderTest(providedValue))
          .rejects
          .toThrow(`unable to coerce file to array; error was ${dataGetErr.message}`)
      })
    })
    describe('nodeApiClient.dataGet doesn\'t err', () => {
      it('should return expected result', async () => {
        /* arrange */
        const providedValue = { file: 'dummyFile' }

        const objectValue = { prop1: 'prop1Value' }
        const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify(objectValue)) }
        nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

        /* act */
        const actualResult = await objectUnderTest(providedValue)

        /* assert */
        expect(actualResult).toEqual({ object: objectValue })
      })
    })
  })
  describe('value is number', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { number: 2.2 }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce number to object; incompatible types`)
    })
  })
  describe('value is object', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { object: { prop1: 'prop1Value' } }

      /* act */
      const actualResult = await objectUnderTest(providedValue)

      /* assert */
      expect(actualResult).toEqual(providedValue)
    })
  })
  describe('value is socket', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { socket: 'dummySocket' }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce socket '${providedValue.socket}' to object; incompatible types`)
    })
  })
  describe('value is string', () => {
    describe('string is JSON', () => {
      describe('JSON is object', () => {
        it('should return expected result', async () => {
          /* arrange */
          const objectValue = { prop1: 'prop1Value' }
          const providedValue = { string: JSON.stringify(objectValue) }

          /* act */
          const actualResult = await objectUnderTest(providedValue)

          /* assert */
          expect(actualResult).toEqual({ object: JSON.parse(providedValue.string) })
        })
      })
      describe('JSON is array', async () => {
        it('returns expected result', async () => {
          /* arrange */
          const providedValue = { string: JSON.stringify([]) }

          /* act/assert */
          await expect(objectUnderTest(providedValue))
            .rejects
            .toThrow(`unable to coerce string to object; error was parsed array but expected object`)
        })
      })
      describe('JSON is null', async () => {
        it('returns expected result', async () => {
          /* arrange */
          const providedValue = { string: JSON.stringify(null) }

          /* act/assert */
          await expect(objectUnderTest(providedValue))
            .rejects
            .toThrow(`unable to coerce string to object; error was parsed null but expected object`)
        })
      })
      describe('JSON is string', async () => {
        it('returns expected result', async () => {
          /* arrange */
          const providedValue = { string: JSON.stringify('dummyString') }

          /* act/assert */
          await expect(objectUnderTest(providedValue))
            .rejects
            .toThrow(`unable to coerce string to object; error was parsed string but expected object`)
        })
      })
    })
    describe('string isn\'t JSON', () => {
      it('returns expected result', async () => {
        /* arrange */
        const providedValue = { string: 'notValidJsonObject' }

        /* act/assert */
        await expect(objectUnderTest(providedValue))
          .rejects
          .toThrow(`unable to coerce string to object; error was Unexpected token o in JSON at position 1`)
      })
    })
  })
  describe('value isnt any of the above', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedValue = { unknown: 'providedValue' }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to object`)
    })
  })
})
