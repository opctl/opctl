jest.mock('./array/validate', () => jest.fn())
import arrayValidate from './array/validate'

jest.mock('./boolean/validate', () => jest.fn())
import booleanValidate from './boolean/validate'

jest.mock('./dir/validate', () => jest.fn())
import dirValidate from './dir/validate'

jest.mock('./file/validate', () => jest.fn())
import fileValidate from './file/validate'

jest.mock('./number/validate', () => jest.fn())
import numberValidate from './number/validate'

jest.mock('./object/validate', () => jest.fn())
import objectValidate from './object/validate'

jest.mock('./socket/validate', () => jest.fn())
import socketValidate from './socket/validate'

jest.mock('./string/validate', () => jest.fn())
import stringValidate from './string/validate'

import objectUnderTest from './validate'

describe('validate', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('param falsy', () => {
    it('returns expected result', () => {
      /* arrange/act/assert */
      expect(() => objectUnderTest({}, null as any))
        .toThrow('param required')
    })
  })
  describe('param.array not falsy', () => {
    it('it calls arrayValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { array: ['dummyItem'] }
      const providedParam = { array: { constraints: { allOf: [] } } }

      const expectedResult = [new Error('dummyError')]
      ; (arrayValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        arrayValidate
      ).toBeCalledWith(
        providedValue.array,
        providedParam.array.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.boolean not falsy', () => {
    it('it calls booleanValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { boolean: true }
      const providedParam = { boolean: { } }

      const expectedResult = [new Error('dummyError')]
      ;(booleanValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        booleanValidate
      ).toBeCalledWith(
        providedValue.boolean
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.dir not falsy', () => {
    it('it calls dirValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { link: 'dummyDir' }
      const providedParam = { dir: { } }

      const expectedResult = [new Error('dummyError')]
      ; (dirValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        dirValidate
      ).toBeCalledWith(
        providedValue.link
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.file not falsy', () => {
    it('it calls fileValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { link: 'dummyFile' }
      const providedParam = { file: { } }

      const expectedResult = [new Error('dummyError')]
      ; (fileValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        fileValidate
      ).toBeCalledWith(
        providedValue.link
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.number not falsy', () => {
    it('it calls numberValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { number: 23 }
      const providedParam = { number: { constraints: { allOf: [ 2, 3 ] } } }

      const expectedResult = [new Error('dummyError')]
      ; (numberValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        numberValidate
      ).toBeCalledWith(
        providedValue.number,
        providedParam.number.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.object not falsy', () => {
    it('it calls objectValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { object: { prop1: 1 } }
      const providedParam = { object: { constraints: { allOf: [ { prop1: 1 } ] } } }

      const expectedResult = [new Error('dummyError')]
      ; (objectValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        objectValidate
      ).toBeCalledWith(
        providedValue.object,
        providedParam.object.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.socket not falsy', () => {
    it('it calls socketValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { socket: 'dummySocket' }
      const providedParam = { socket: { } }

      const expectedResult = [new Error('dummyError')]
      ; (socketValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        socketValidate
      ).toBeCalledWith(
        providedValue.socket
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.string not falsy', () => {
    it('it calls stringValidate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { string: 'dummyString' }
      const providedParam = { string: { constraints: { format: '.*' } } }

      const expectedResult = [new Error('dummyError')]
      ; (stringValidate as jest.Mock).mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        stringValidate
      ).toBeCalledWith(
        providedValue.string,
        providedParam.string.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
})
