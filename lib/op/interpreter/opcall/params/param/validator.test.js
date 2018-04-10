const arrayValidator = require('./array/validator')
const booleanValidator = require('./boolean/validator')
const dirValidator = require('./dir/validator')
const fileValidator = require('./file/validator')
const numberValidator = require('./number/validator')
const objectValidator = require('./object/validator')
const socketValidator = require('./socket/validator')
const stringValidator = require('./string/validator')

const objectUnderTest = require('./validator')

describe('validate', () => {
  beforeEach(() => {
    jest.resetAllMocks()
  })
  describe('param falsy', () => {
    it('returns expected result', () => {
      /* arrange/act/assert */
      expect(() => objectUnderTest.validate({}, null))
        .toThrow('param required')
    })
  })
  describe('param.array not falsy', () => {
    it('it calls arrayValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { array: ['dummyItem'] }
      const providedParam = { array: { constraints: { allOf: [] } } }

      const expectedResult = [new Error('dummyError')]
      arrayValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        arrayValidator.validate
      ).toBeCalledWith(
        providedValue.array,
        providedParam.array.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.boolean not falsy', () => {
    it('it calls booleanValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { boolean: true }
      const providedParam = { boolean: { } }

      const expectedResult = [new Error('dummyError')]
      booleanValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        booleanValidator.validate
      ).toBeCalledWith(
        providedValue.boolean
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.dir not falsy', () => {
    it('it calls dirValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { dir: 'dummyDir' }
      const providedParam = { dir: { } }

      const expectedResult = [new Error('dummyError')]
      dirValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        dirValidator.validate
      ).toBeCalledWith(
        providedValue.dir
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.file not falsy', () => {
    it('it calls fileValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { file: 'dummyFile' }
      const providedParam = { file: { } }

      const expectedResult = [new Error('dummyError')]
      fileValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        fileValidator.validate
      ).toBeCalledWith(
        providedValue.file
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.number not falsy', () => {
    it('it calls numberValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { number: 23 }
      const providedParam = { number: { constraints: { allOf: [ 2, 3 ] } } }

      const expectedResult = [new Error('dummyError')]
      numberValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        numberValidator.validate
      ).toBeCalledWith(
        providedValue.number,
        providedParam.number.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.object not falsy', () => {
    it('it calls objectValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { object: { prop1: 1 } }
      const providedParam = { object: { constraints: { allOf: [ { prop1: 1 } ] } } }

      const expectedResult = [new Error('dummyError')]
      objectValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        objectValidator.validate
      ).toBeCalledWith(
        providedValue.object,
        providedParam.object.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.socket not falsy', () => {
    it('it calls socketValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { socket: 'dummySocket' }
      const providedParam = { socket: { } }

      const expectedResult = [new Error('dummyError')]
      socketValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        socketValidator.validate
      ).toBeCalledWith(
        providedValue.socket
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('param.string not falsy', () => {
    it('it calls stringValidator.validate w/ expected args & returns result', () => {
      /* arrange */
      const providedValue = { string: 'dummyString' }
      const providedParam = { string: { constraints: { format: '.*' } } }

      const expectedResult = [new Error('dummyError')]
      stringValidator.validate = jest.fn()
        .mockReturnValue(expectedResult)

      /* act */
      const actualResult = objectUnderTest.validate(
        providedValue,
        providedParam
      )

      /* assert */
      expect(
        stringValidator.validate
      ).toBeCalledWith(
        providedValue.string,
        providedParam.string.constraints
      )

      expect(actualResult).toEqual(expectedResult)
    })
  })
})
