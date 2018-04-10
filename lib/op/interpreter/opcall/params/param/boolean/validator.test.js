const objectUnderTest = require('./validator')

describe('validate', () => {
  describe('typeof value isn\'t boolean', () => {
    it('returns expected result', () => {
      /* arrange */
      const expectedResult = [new Error('boolean required')]

      /* act */
      const actualResult = objectUnderTest.validate(undefined)

      /* assert */
      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('typeof value is boolean', () => {
    it('returns expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest.validate(true)

      /* assert */
      expect(actualResult).toEqual([])
    })
  })
})
