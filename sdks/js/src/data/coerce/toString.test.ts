jest.mock('../../api/client/data/get', () => jest.fn())
import dataGet from '../../api/client/data/get'

import objectUnderTest from './toString'
import Value from '../../types/value';

describe('toString', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('value is array', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { array: ['item1'] }

      /* act */
      const actualResult = await objectUnderTest(
        'apiBaseUrl',
        providedValue
      )

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.array) })
    })
  })
  describe('value is boolean', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { boolean: false }

      /* act */
      const actualResult = await objectUnderTest(
        'apiBaseUrl',
        providedValue
      )

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.boolean) })
    })
  })
  describe('value is file', () => {
    it('should call dataGet w/ expected args', async () => {
      /* arrange */
      const providedNetRef = 'providedNetRef'
      const providedValue = { file: 'dummyFile' }


      const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify({})) }
        ; (dataGet as jest.Mock).mockResolvedValue(dataGetResult)

      /* act */
      await objectUnderTest(
        providedNetRef,
        providedValue
      )

      /* assert */
      expect(dataGet)
        .toBeCalledWith(
          providedNetRef,
          providedValue.file
        )
    })
    describe('dataGet errs', () => {
      it('returns expected result', async () => {
        /* arrange */
        const providedValue = { file: 'dummyFile' }

        const dataGetErr = new Error('dummyErr')
          ; (dataGet as jest.Mock).mockRejectedValue(dataGetErr)

        /* act/assert */
        await expect(objectUnderTest(
          'apiBaseUrl',
          providedValue
        ))
          .rejects
          .toThrow(`unable to coerce file to string; error was ${dataGetErr.message}`)
      })
    })
    describe('dataGet doesn\'t err', () => {
      it('should return expected result', async () => {
        /* arrange */
        const providedValue = { file: 'dummyFile' }

        const textResult = 'textResult'
        const dataGetResult = { text: jest.fn().mockResolvedValue(textResult) }
          ; (dataGet as jest.Mock).mockResolvedValue(dataGetResult)

        /* act */
        const actualResult = await objectUnderTest(
          'apiBaseUrl',
          providedValue
        )

        /* assert */
        expect(actualResult).toEqual({ string: textResult })
      })
    })
  })
  describe('value is number', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { number: 2.2 }

      /* act */
      const actualResult = await objectUnderTest(
        'apiBaseUrl',
        providedValue
      )

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.number) })
    })
  })
  describe('value is object', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { object: { prop1: 'value1' } }

      /* act */
      const actualResult = await objectUnderTest(
        'apiBaseUrl',
        providedValue
      )

      /* assert */
      expect(actualResult).toEqual({ string: JSON.stringify(providedValue.object) })
    })
  })
  describe('value is string', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { string: 'value' }

      /* act */
      const actualResult = await objectUnderTest(
        'apiBaseUrl',
        providedValue
      )

      /* assert */
      expect(actualResult).toEqual(providedValue)
    })
  })
  describe('value isnt any of the above', () => {
    it('returns expected result', async () => {
      /* arrange */
      const providedValue = { unknown: 'providedValue' } as Value

      /* act/assert */
      await expect(objectUnderTest(
        'apiBaseUrl',
        providedValue
      ))
        .rejects
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to string`)
    })
  })
})
