const objectUnderTest = require('./validator')

describe('validate', () => {
  describe('constraints undefined', () => {
    test('returns expected result', () => {
      /* arrange/act */
      const actualErrors = objectUnderTest.validate(
        'dummyValue'
      )

      /* assert */
      expect(actualErrors).toEqual([])
    })
  })
  describe('constraints defined', () => {
    describe('format constraint', () => {
      describe('constraint.format is docker-image-ref', () => {
        describe('value is docker-image-ref', () => {
          test('returns expected result', () => {
            /* arrange/act */
            const actualErrors = objectUnderTest.validate('owner/registry:tag', {
              format: 'docker-image-ref'
            })

            /* assert */
            expect(actualErrors).toEqual([])
          })
        })
        describe('value isnt docker-image-ref', () => {
          test('returns expected result', () => {
            /* arrange/act */
            const actualErrors = objectUnderTest.validate('#', {
              format: 'docker-image-ref'
            })

            /* assert */
            expect(actualErrors).toEqual([{
              'dataPath': '',
              'keyword': 'format',
              'message': 'should match format "docker-image-ref"',
              'params': {'format': 'docker-image-ref'},
              'schemaPath': '#/format'
            }])
          })
        })
      })
    })
    describe('minLength constraint', () => {
      describe('value > minLength', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            'dummyValue',
            {
              minLength: 1
            }
          )

          /* assert */
          expect(actualErrors).toEqual([])
        })
      })
      describe('value < minLength', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            'dummyValue',
            {
              minLength: 100
            }
          )

          /* assert */
          expect(actualErrors).toEqual([{
            'dataPath': '',
            'keyword': 'minLength',
            'message': 'should NOT be shorter than 100 characters',
            'params': {'limit': 100},
            'schemaPath': '#/minLength'
          }])
        })
      })
    })
  })
})
