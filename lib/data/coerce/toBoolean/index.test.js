const nodeApiClient = require('../../../node/api/client')
jest.mock('./isStringTruthy', () => jest.fn())
const isStringTruthy = require('./isStringTruthy')

const objectUnderTest = require('./index')

describe('toBoolean', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('value is array', () => {
    describe('value is empty', () => {
      it('should return expected result', async () => {
        /* arrange */
        const providedValue = { array: [] }

        /* act */
        const actualResult = await objectUnderTest.toBoolean(providedValue)

        /* assert */
        expect(actualResult).toEqual({ boolean: false })
      })
    })
    describe('value isn\'t empty', () => {
      it('should return expected result', async () => {
        /* arrange */
        const providedValue = { array: ['item1'] }

        /* act */
        const actualResult = await objectUnderTest.toBoolean(providedValue)

        /* assert */
        expect(actualResult).toEqual({ boolean: true })
      })
    })
  })
  describe('value is boolean', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { boolean: false }

      /* act */
      const actualResult = await objectUnderTest.toBoolean(providedValue)

      /* assert */
      expect(actualResult).toEqual({ boolean: false })
    })
  })
  describe('value is dir', () => {
    it('should call nodeApiClient.dataGet w/ expected args', async () => {
      /* arrange */
      const providedValue = { dir: 'dummyDir' }

      const dataGetResult = { json: () => ([]) }
      nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

      /* act */
      await objectUnderTest.toBoolean(providedValue)

      /* assert */
      expect(nodeApiClient.dataGet)
        .toBeCalledWith({ dataRef: providedValue.dir })
    })
    describe('nodeApiClient.dataGet errs', () => {
      it('returns expected result', async () => {
        /* arrange */
        const providedValue = { dir: 'dummyFile' }

        const dataGetErr = new Error('dummyErr')
        nodeApiClient.dataGet = jest.fn().mockRejectedValue(dataGetErr)

        /* act/assert */
        await expect(objectUnderTest.toBoolean(providedValue))
          .rejects
          .toThrow(`unable to coerce dir to boolean; error was ${dataGetErr.message}`)
      })
    })
    describe('nodeApiClient.dataGet doesn\'t err', () => {
      describe('array is empty', () => {
        it('should return expected result', async () => {
          /* arrange */
          const providedValue = { dir: 'dummyDir' }

          const dataGetResult = { json: jest.fn().mockResolvedValue([]) }
          nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

          /* act */
          const actualResult = await objectUnderTest.toBoolean(providedValue)

          /* assert */
          expect(actualResult).toEqual({ boolean: false })
        })
      })
      describe('array isn\'t empty', () => {
        it('should return expected result', async () => {
          /* arrange */
          const providedValue = { dir: 'dummyDir' }

          const dataGetResult = { json: jest.fn().mockResolvedValue(['item']) }
          nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

          /* act */
          const actualResult = await objectUnderTest.toBoolean(providedValue)

          /* assert */
          expect(actualResult).toEqual({ boolean: true })
        })
      })
    })
  })
  describe('value is file', () => {
    it('should call nodeApiClient.dataGet w/ expected args', async () => {
      /* arrange */
      const providedValue = { file: 'dummyFile' }

      const dataGetResult = { text: () => ('dummyText') }
      nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

      /* act */
      await objectUnderTest.toBoolean(providedValue)

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
        await expect(objectUnderTest.toBoolean(providedValue))
          .rejects
          .toThrow(`unable to coerce file to boolean; error was ${dataGetErr.message}`)
      })
    })
    describe('nodeApiClient.dataGet doesn\'t err', () => {
      describe('should call isStringTruthy w/ expected args & return expected result', () => {
        it('should return expected result', async () => {
          /* arrange */
          const providedValue = { file: 'dummyFile' }

          const textResult = 'textResult'
          const dataGetResult = { text: jest.fn().mockResolvedValue(textResult) }
          nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

          const isStringTruthyResult = true
          isStringTruthy.mockImplementation(() => isStringTruthyResult)

          /* act */
          const actualResult = await objectUnderTest.toBoolean(providedValue)

          /* assert */
          expect(isStringTruthy)
            .toBeCalledWith(textResult)

          expect(actualResult).toEqual({ boolean: isStringTruthyResult })
        })
      })
    })
  })
  describe('value is number', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { number: 2.2 }

      /* act */
      const actualResult = await objectUnderTest.toBoolean(providedValue)

      /* assert */
      expect(actualResult).toEqual({ boolean: true })
    })
  })
  describe('value is object', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { object: { prop1: 'value1' } }

      /* act */
      const actualResult = await objectUnderTest.toBoolean(providedValue)

      /* assert */
      expect(actualResult).toEqual({ boolean: true })
    })
  })
  describe('value is socket', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedValue = { socket: 'providedValue' }

      /* act/assert */
      await expect(objectUnderTest.toBoolean(providedValue))
        .rejects
        .toThrow(`unable to coerce socket '${providedValue.socket}' to boolean; incompatible types`)
    })
  })
  describe('value is string', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { string: 'value' }

      const isStringTruthyResult = true
      isStringTruthy.mockImplementation(() => isStringTruthyResult)

      /* act */
      const actualResult = await objectUnderTest.toBoolean(providedValue)

      /* assert */
      expect(isStringTruthy)
        .toBeCalledWith(providedValue.string)

      expect(actualResult).toEqual({ boolean: isStringTruthyResult })
    })
  })
  describe('value isnt any of the above', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedValue = { unknown: 'providedValue' }

      /* act/assert */
      await expect(objectUnderTest.toBoolean(providedValue))
        .rejects
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to boolean`)
    })
  })
})
