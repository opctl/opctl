import objectUnderTest from './validate'

describe('validate', () => {
  describe('constraints undefined', () => {
    it('returns expected result', () => {
      /* arrange/act */
      const actualErrors = objectUnderTest(
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
          it('returns expected result', () => {
            /* arrange/act */
            const actualErrors = objectUnderTest(1, {
              format: 'integer'
            })

            /* assert */
            expect(actualErrors).toEqual([])
          })
        })
        describe('value isnt integer', () => {
          it('returns expected result', () => {
            /* arrange/act */
            const actualErrors = objectUnderTest(1.1, {
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
        it('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest(
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
        it('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest(
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
