const objectUnderTest = require('./validator')

describe('validate', () => {
  describe('constraints undefined', () => {
    it('returns expected result', () => {
      /* arrange/act */
      const actualErrors = objectUnderTest.validate(
        []
      )

      /* assert */
      expect(actualErrors).toEqual([])
    })
  })
  describe('constraints defined', () => {
    describe('minItems constraint', () => {
      describe('value.length > minItems', () => {
        it('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            ['dummyItem'],
            {
              minItems: 1
            }
          )

          /* assert */
          expect(actualErrors).toEqual([])
        })
      })
      describe('value.length < minItems', () => {
        it('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            [],
            {
              minItems: 1
            }
          )

          /* assert */
          expect(actualErrors).toEqual([{
            'dataPath': '',
            'keyword': 'minItems',
            'message': 'should NOT have less than 1 items',
            'params': {'limit': 1},
            'schemaPath': '#/minItems'
          }])
        })
      })
    })
  })
})
