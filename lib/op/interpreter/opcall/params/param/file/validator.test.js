const objectUnderTest = require('./validator')

describe('validate', () => {
  describe('value falsy', () => {
    it('returns expected result', () => {
      /* arrange */
      const expectedResult = [new Error('file required')]

      /* act */
      const actualResult = objectUnderTest.validate('')

      /* assert */
      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('value truthy', () => {
    it('returns expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest.validate('some file')

      /* assert */
      expect(actualResult).toEqual([])
    })
  })
})
