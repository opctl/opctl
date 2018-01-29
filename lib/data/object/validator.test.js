const objectUnderTest = require('./validator');

describe('validate', () => {
  describe('constraints undefined', () => {
    test('returns expected result', () => {
      /* arrange/act */
      const actualErrors = objectUnderTest.validate(
        {}
      );

      /* assert */
      expect(actualErrors).toEqual([])
    })
  });
  describe('constraints defined', () => {
    describe('minProperties constraint', () => {
      describe('Object.keys(value).length > minProperties', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            {'dummyProperty': 'dummyValue'},
            {
              minProperties: 1
            }
          );

          /* assert */
          expect(actualErrors).toEqual([])
        })
      });
      describe('Object.keys(value).length < minProperties', () => {
        test('returns expected result', () => {
          /* arrange/act */
          const actualErrors = objectUnderTest.validate(
            {},
            {
              minProperties: 1
            }
          );

          /* assert */
          expect(actualErrors).toEqual([{
            "dataPath": "",
            "keyword": "minProperties",
            "message": "should NOT have less than 1 properties",
            "params": {"limit": 1},
            "schemaPath": "#/minProperties"
          }])
        })
      });
    })
  })
});
