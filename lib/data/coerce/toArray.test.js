const nodeApiClient = require('../../node/api/client')

const objectUnderTest = require('./toArray')

describe('toArray', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('value is array', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { array: ['item1'] }

      /* act */
      const actualResult = await objectUnderTest(providedValue)

      /* assert */
      expect(actualResult).toEqual(providedValue)
    })
  })
  describe('value is boolean', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { boolean: true }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce boolean '${providedValue.boolean}' to array; incompatible types`)
    })
  })
  describe('value is dir', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { dir: 'dummyDir' }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce dir '${providedValue.dir}' to array; incompatible types`)
    })
  })
  describe('value is file', () => {
    it('should call nodeApiClient.dataGet w/ expected args', async () => {
      /* arrange */
      const providedValue = { file: 'dummyFile' }

      const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify([])) }
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

        const arrayValue = ['item1']
        const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify(arrayValue)) }
        nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

        /* act */
        const actualResult = await objectUnderTest(providedValue)

        /* assert */
        expect(actualResult).toEqual({ array: arrayValue })
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
        .toThrow(`unable to coerce number to array; incompatible types`)
    })
  })
  describe('value is socket', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { socket: 'dummySocket' }

      /* act/assert */
      await expect(objectUnderTest(providedValue))
        .rejects
        .toThrow(`unable to coerce socket '${providedValue.socket}' to array; incompatible types`)
    })
  })
  describe('value is string', () => {
    describe('string is JSON', () => {
      describe('JSON is array', () => {
        it('should return expected result', async () => {
          /* arrange */
          const arrayValue = ['item1']
          const providedValue = { string: JSON.stringify(arrayValue) }

          /* act */
          const actualResult = await objectUnderTest(providedValue)

          /* assert */
          expect(actualResult).toEqual({ array: JSON.parse(providedValue.string) })
        })
      })
      describe('JSON isn\'t array', async () => {
        it('returns expected result', async () => {
          /* arrange */
          const objectValue = { prop1: 'prop1Value' }
          const providedValue = { string: JSON.stringify(objectValue) }

          /* act/assert */
          await expect(objectUnderTest(providedValue))
            .rejects
            .toThrow(`unable to coerce string to array; error was parsed ${typeof objectValue} but expected array`)
        })
      })
    })
    describe('string isn\'t JSON', () => {
      it('returns expected result', async () => {
        /* arrange */
        const providedValue = { string: 'notValidJSONArray' }

        /* act/assert */
        await expect(objectUnderTest(providedValue))
          .rejects
          .toThrow(`unable to coerce string to array; error was Unexpected token o in JSON at position 1`)
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
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to array`)
    })
  })
})
