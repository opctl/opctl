import objectUnderTest from './validate'

describe('validate', () => {
  describe('typeof value isn\'t boolean', () => {
    it('returns expected result', () => {
      /* arrange */
      const expectedResult = [new Error('boolean required')]

      /* act */
      const actualResult = objectUnderTest(
        undefined as unknown as boolean
      )

      /* assert */
      expect(actualResult).toEqual(expectedResult)
    })
  })
  describe('typeof value is boolean', () => {
    it('returns expected result', () => {
      /* arrange/act */
      const actualResult = objectUnderTest(true)

      /* assert */
      expect(actualResult).toEqual([])
    })
  })
})
