jest.mock('../../api/client/data/get', () => jest.fn())
import dataGet from '../../api/client/data/get'

import objectUnderTest from './toNumber'
import Value from '../../model/value'

describe('toNumber', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('value is array', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { array: [] }

      /* act/assert */
      await expect(objectUnderTest(
        'apiBaseUrl',
        providedValue
      ))
        .rejects
        .toThrow(`unable to coerce array to number; incompatible types`)
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
        .toThrow(`unable to coerce boolean '${providedValue.boolean}' to number; incompatible types`)
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
        .toThrow(`unable to coerce dir '${providedValue.dir}' to number; incompatible types`)
    })
  })
  describe('value is file', () => {
    it('should call dataGet w/ expected args', async () => {
      /* arrange */
      const providedNetRef = 'providedNetRef'
      const providedValue = { file: 'dummyFile' }

      const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify(2.2)) }
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
          .toThrow(`unable to coerce file to number; error was ${dataGetErr.message}`)
      })
    })
    describe('dataGet doesn\'t err', () => {
      it('should return expected result', async () => {
        /* arrange */
        const providedValue = { file: 'dummyFile' }

        const numberValue = 3.2
        const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify(numberValue)) }
          ; (dataGet as jest.Mock).mockResolvedValue(dataGetResult)

        /* act */
        const actualResult = await objectUnderTest(
          'apiBaseUrl',
          providedValue
        )

        /* assert */
        expect(actualResult).toEqual({ number: numberValue })
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
      expect(actualResult).toEqual(providedValue)
    })
  })
  describe('value is object', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { object: {} }

      /* act/assert */
      await expect(objectUnderTest(
        'apiBaseUrl',
        providedValue
      ))
        .rejects
        .toThrow(`unable to coerce object to number; incompatible types`)
    })
  })
  describe('value is socket', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { socket: 'dummySocket' }

      /* act/assert */
      await expect(objectUnderTest(
        'apiBaseUrl',
        providedValue
      ))
        .rejects
        .toThrow(`unable to coerce socket '${providedValue.socket}' to number; incompatible types`)
    })
  })
  describe('value is string', () => {
    describe('string is JSON', () => {
      describe('JSON is number', () => {
        it('should return expected result', async () => {
          /* arrange */
          const numberValue = 3.3
          const providedValue = { string: JSON.stringify(numberValue) }

          /* act */
          const actualResult = await objectUnderTest(
            'apiBaseUrl',
            providedValue
          )

          /* assert */
          expect(actualResult).toEqual({ number: JSON.parse(providedValue.string) })
        })
      })
      describe('JSON isn\'t number', () => {
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
            .toThrow(`unable to coerce string to number; error was parsed ${typeof objectValue} but expected number`)
        })
      })
    })
    describe('string isn\'t JSON', () => {
      it('returns expected result', async () => {
        /* arrange */
        const providedValue = { string: 'notValidJsonNumber' }

        /* act/assert */
        await expect(objectUnderTest(
          'apiBaseUrl',
          providedValue
        ))
          .rejects
          .toThrow(`unable to coerce string to number; error was Unexpected token 'o', \"notValidJsonNumber\" is not valid JSON`)
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
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to number`)
    })
  })
})
