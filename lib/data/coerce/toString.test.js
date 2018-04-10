const nodeApiClient = require('../../node/api/client')

const objectUnderTest = require('./toString')

describe('toString', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('value is array', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { array: ['item1'] }

      /* act */
      const actualResult = await objectUnderTest.toString(providedValue)

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.array) })
    })
  })
  describe('value is boolean', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { boolean: false }

      /* act */
      const actualResult = await objectUnderTest.toString(providedValue)

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.boolean) })
    })
  })
  describe('value is file', () => {
    it('should call nodeApiClient.dataGet w/ expected args', async () => {
      /* arrange */
      const providedValue = { file: 'dummyFile' }

      nodeApiClient.dataGet = jest.fn().mockResolvedValue({ text () { } })

      /* act */
      await objectUnderTest.toString(providedValue)

      /* assert */
      expect(nodeApiClient.dataGet)
        .toBeCalledWith({ dataRef: providedValue.file })
    })
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { file: 'dummyFile' }

      const textResult = 'textResult'
      const dataGetResult = { text: jest.fn().mockResolvedValue(textResult) }
      nodeApiClient.dataGet = jest.fn().mockResolvedValue(dataGetResult)

      /* act */
      const actualResult = await objectUnderTest.toString(providedValue)

      /* assert */
      expect(actualResult).toEqual({ string: textResult })
    })
  })
  describe('value is number', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { number: 2.2 }

      /* act */
      const actualResult = await objectUnderTest.toString(providedValue)

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.number) })
    })
  })
  describe('value is object', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { object: { prop1: 'value1' } }

      /* act */
      const actualResult = await objectUnderTest.toString(providedValue)

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.object) })
    })
  })
  describe('value is string', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { string: 'value' }

      /* act */
      const actualResult = await objectUnderTest.toString(providedValue)

      /* assert */
      expect(actualResult).toEqual(providedValue)
    })
  })
  describe('value isnt array/boolean/file/number/object/string', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedValue = { dir: 'providedValue' }
      /* act/assert */

      await expect(objectUnderTest.toString(providedValue))
        .rejects
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to string`)
    })
  })
})
