import objectUnderTest from './validate'

describe('validate', () => {
  describe('constraints undefined', () => {
    it('returns expected result', () => {
      /* arrange/act */
      const actualErrors = objectUnderTest(
        {}
      )

      /* assert */
      expect(actualErrors).toEqual([])
    })
  })
  describe('constraints defined', () => {
    describe('minProperties constraint', () => {
      describe('Object.keys(value).length > minProperties', () => {
        it('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest(
            {'dummyProperty': 'dummyValue'},
            {
              minProperties: 1
            }
          )

          /* assert */
          expect(actualErrors).toEqual([])
        })
      })
      describe('Object.keys(value).length < minProperties', () => {
        it('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest(
            {},
            {
              minProperties: 1
            }
          )

          /* assert */
          expect(actualErrors).toEqual([{
            'dataPath': '',
            'keyword': 'minProperties',
            'message': 'should NOT have less than 1 properties',
            'params': {'limit': 1},
            'schemaPath': '#/minProperties'
          }])
        })
      })
    })
  })
})
