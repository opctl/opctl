import objectUnderTest from './isStringTruthy'

describe('isStringTruthy', () => {
  describe('string is all zeros', () => {
    it('should return expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest('00')

      /* assert */
      expect(actualResult).toEqual(false)
    })
  })
  describe('string is empty', () => {
    it('should return expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest('')

      /* assert */
      expect(actualResult).toEqual(false)
    })
  })
  describe('string is \'F\'', () => {
    it('should return expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest('F')

      /* assert */
      expect(actualResult).toEqual(false)
    })
  })
  describe('string is \'FALSE\'', () => {
    it('should return expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest('FALSE')

      /* assert */
      expect(actualResult).toEqual(false)
    })
  })
  describe('string is none of the above', () => {
    it('should return expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest('noneOfTheAbove')

      /* assert */
      expect(actualResult).toEqual(true)
    })
  })
})
