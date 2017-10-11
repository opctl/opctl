const objectUnderTest = require('./validator');

describe('validate', () => {
  describe('constraints undefined', () => {
    test('returns expected result', () => {
      /* arrange/act */
      const actualErrors = objectUnderTest.validate(
        1
      );

      /* assert */
      expect(actualErrors).toEqual([])
    })
  });
  describe('constraints defined', () => {
    describe('minimum constraint', () => {
      describe('value > minimum', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            2,
            {
              minimum: 1
            }
          );

          /* assert */
          expect(actualErrors).toBeNull()
        })
      });
      describe('value < minimum', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            1,
            {
              minimum: 2
            }
          );

          /* assert */
          expect(actualErrors).toEqual([{
            "dataPath": "",
            "keyword": "minimum",
            "message": "should be >= 2",
            "params": {"comparison": ">=", "exclusive": false, "limit": 2},
            "schemaPath": "#/minimum"
          }])
        })
      });
    })
  })
});
