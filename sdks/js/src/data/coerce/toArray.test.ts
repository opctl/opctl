jest.mock('../../api/client/data/get', () => jest.fn())
import dataGet from '../../api/client/data/get'

import objectUnderTest from './toArray'
import Value from '../../types/value'

describe('toArray', () => {
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
      expect(actualResult).toEqual(providedValue)
    })
  })
  describe('value is boolean', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { boolean: true }

      /* act/assert */
      await expect(objectUnderTest(
        'apiBaseUrl',
        providedValue
      ))
        .rejects
        .toThrow(`unable to coerce boolean '${providedValue.boolean}' to array; incompatible types`)
    })
  })
  describe('value is dir', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { dir: 'dummyDir' }

      /* act/assert */
      await expect(objectUnderTest(
        'apiBaseUrl',
        providedValue
      ))
        .rejects
        .toThrow(`unable to coerce dir '${providedValue.dir}' to array; incompatible types`)
    })
  })
  describe('value is file', () => {
    it('should call dataGet w/ expected args', async () => {
      /* arrange */
      const providedNetRef = 'providedNetRef'
      const providedValue = { file: 'dummyFile' }

      const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify([])) }
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
        await expect(objectUnderTest('apiBaseUrl', providedValue))
          .rejects
          .toThrow(`unable to coerce file to array; error was ${dataGetErr.message}`)
      })
    })
    describe('dataGet doesn\'t err', () => {
      it('should return expected result', async () => {
        /* arrange */
        const providedValue = { file: 'dummyFile' }

        const arrayValue = ['item1']
        const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify(arrayValue)) }
          ; (dataGet as jest.Mock).mockResolvedValue(dataGetResult)

        /* act */
        const actualResult = await objectUnderTest(
          'apiBaseUrl',
          providedValue
        )

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
      await expect(objectUnderTest('apiBaseUrl', providedValue))
        .rejects
        .toThrow(`unable to coerce number to array; incompatible types`)
    })
  })
  describe('value is object', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { object: {} }

      /* act/assert */
      await expect(objectUnderTest('apiBaseUrl', providedValue))
        .rejects
        .toThrow(`unable to coerce object to array; incompatible types`)
    })
  })
  describe('value is socket', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { socket: 'dummySocket' }

      /* act/assert */
      await expect(objectUnderTest('apiBaseUrl', providedValue))
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
          const actualResult = await objectUnderTest(
            'apiBaseUrl',
            providedValue
          )

          /* assert */
          expect(actualResult).toEqual({ array: JSON.parse(providedValue.string) })
        })
      })
      describe('JSON isn\'t array', () => {
        it('returns expected result', async () => {
          /* arrange */
          const objectValue = { prop1: 'prop1Value' }
          const providedValue = { string: JSON.stringify(objectValue) }

          /* act/assert */
          await expect(objectUnderTest(
            'apiBaseUrl',
            providedValue
          ))
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
        await expect(objectUnderTest(
          'apiBaseUrl',
          providedValue
        ))
          .rejects
          .toThrow(`unable to coerce string to array; error was Unexpected token o in JSON at position 1`)
      })
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
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to array`)
    })
  })
})
