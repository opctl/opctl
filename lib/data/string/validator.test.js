const objectUnderTest = require('./validator');

describe('validate', () => {
  describe('constraints undefined', () => {
    test('throws', () => {
      /* arrange */
      const call = () => objectUnderTest.validate('dummyValue');

      /* act/assert */
      expect(call).toThrow();
    })
  });
  describe('constraints defined', () => {
    describe('minLength constraint', () => {
      describe('value > minLength', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            'dummyValue',
            {
              minLength: 1
            }
          );

          /* assert */
          expect(actualErrors).toBeNull()
        })
      });
      describe('value < minLength', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            'dummyValue',
            {
              minLength: 100
            }
          );

          /* assert */
          expect(actualErrors).toEqual([{
            "dataPath": "",
            "keyword": "minLength",
            "message": "should NOT be shorter than 100 characters",
            "params": {"limit": 100},
            "schemaPath": "#/minLength"
          }])
        })
      });
    })
  })
});
