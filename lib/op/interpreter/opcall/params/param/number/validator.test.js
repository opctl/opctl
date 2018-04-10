const objectUnderTest = require('./validator')

describe('validate', () => {
  describe('constraints undefined', () => {
    test('returns expected result', () => {
      /* arrange/act */
      const actualErrors = objectUnderTest.validate(
        1
      )

      /* assert */
      expect(actualErrors).toEqual([])
    })
  })
  describe('constraints defined', () => {
    describe('format constraint', () => {
      describe('constraint.format is integer', () => {
        describe('value is integer', () => {
          test('returns expected result', () => {
            /* arrange/act */
            const actualErrors = objectUnderTest.validate(1, {
              format: 'integer'
            })

            /* assert */
            expect(actualErrors).toEqual([])
          })
        })
        describe('value isnt integer', () => {
          test('returns expected result', () => {
            /* arrange/act */
            const actualErrors = objectUnderTest.validate(1.1, {
              format: 'integer'
            })

            /* assert */
            expect(actualErrors).toEqual([{
              'dataPath': '',
              'keyword': 'format',
              'message': 'should match format "integer"',
              'params': {'format': 'integer'},
              'schemaPath': '#/format'
            }])
          })
        })
      })
    })
    describe('minimum constraint', () => {
      describe('value > minimum', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            2,
            {
              minimum: 1
            }
          )

          /* assert */
          expect(actualErrors).toEqual([])
        })
      })
      describe('value < minimum', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            1,
            {
              minimum: 2
            }
          )

          /* assert */
          expect(actualErrors).toEqual([{
            'dataPath': '',
            'keyword': 'minimum',
            'message': 'should be >= 2',
            'params': {'comparison': '>=', 'exclusive': false, 'limit': 2},
            'schemaPath': '#/minimum'
          }])
        })
      })
    })
  })
})
