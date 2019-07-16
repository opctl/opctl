jest.mock('../../api/client/data/get', () => jest.fn())
import dataGet from '../../api/client/data/get'

import objectUnderTest from './toObject'
import Value from '../../model/value'

describe('toObject', () => {
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
        .toThrow(`unable to coerce array to object; incompatible types`)
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
        .toThrow(`unable to coerce boolean '${providedValue.boolean}' to object; incompatible types`)
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
        .toThrow(`unable to coerce dir '${providedValue.dir}' to object; incompatible types`)
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
          .toThrow(`unable to coerce file to array; error was ${dataGetErr.message}`)
      })
    })
    describe('dataGet doesn\'t err', () => {
      it('should return expected result', async () => {
        /* arrange */
        const providedValue = { file: 'dummyFile' }

        const objectValue = { prop1: 'prop1Value' }
        const dataGetResult = { text: jest.fn().mockResolvedValue(JSON.stringify(objectValue)) }
          ; (dataGet as jest.Mock).mockResolvedValue(dataGetResult)

        /* act */
        const actualResult = await objectUnderTest(
          'apiBaseUrl',
          providedValue
        )

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
      await expect(objectUnderTest(
        'apiBaseUrl',
        providedValue
      ))
        .rejects
        .toThrow(`unable to coerce number to object; incompatible types`)
    })
  })
  describe('value is object', () => {
    it('should return expected result', async () => {
      /* arrange */
      const providedValue = { object: { prop1: 'prop1Value' } }

      /* act */
      const actualResult = await objectUnderTest(
        'apiBaseUrl',
        providedValue
      )

      /* assert */
      expect(actualResult).toEqual(providedValue)
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
          const actualResult = await objectUnderTest(
            'apiBaseUrl',
            providedValue
          )

          /* assert */
          expect(actualResult).toEqual({ object: JSON.parse(providedValue.string) })
        })
      })
      describe('JSON is array', () => {
        it('returns expected result', async () => {
          /* arrange */
          const providedValue = { string: JSON.stringify([]) }

          /* act/assert */
          await expect(objectUnderTest(
            'apiBaseUrl',
            providedValue
          ))
            .rejects
            .toThrow(`unable to coerce string to object; error was parsed array but expected object`)
        })
      })
      describe('JSON is null', () => {
        it('returns expected result', async () => {
          /* arrange */
          const providedValue = { string: JSON.stringify(null) }

          /* act/assert */
          await expect(objectUnderTest(
            'apiBaseUrl',
            providedValue
          ))
            .rejects
            .toThrow(`unable to coerce string to object; error was parsed null but expected object`)
        })
      })
      describe('JSON is string', () => {
        it('returns expected result', async () => {
          /* arrange */
          const providedValue = { string: JSON.stringify('dummyString') }

          /* act/assert */
          await expect(objectUnderTest(
            'apiBaseUrl',
            providedValue
          ))
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
        await expect(objectUnderTest(
          'apiBaseUrl',
          providedValue
        ))
          .rejects
          .toThrow(`unable to coerce string to object; error was Unexpected token o in JSON at position 1`)
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
        .toThrow(`unable to coerce ${JSON.stringify(providedValue)} to object`)
    })
  })
})
