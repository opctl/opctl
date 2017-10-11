'use strict';

var objectUnderTest = require('./validator');

describe('validate', function () {
  describe('constraints undefined', function () {
    test('returns expected result', function () {
      /* arrange/act */
      var actualErrors = objectUnderTest.validate(1);

      /* assert */
      expect(actualErrors).toEqual([]);
    });
  });
  describe('constraints defined', function () {
    describe('minimum constraint', function () {
      describe('value > minimum', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate(2, {
            minimum: 1
          });

          /* assert */
          expect(actualErrors).toBeNull();
        });
      });
      describe('value < minimum', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate(1, {
            minimum: 2
          });

          /* assert */
          expect(actualErrors).toEqual([{
            "dataPath": "",
            "keyword": "minimum",
            "message": "should be >= 2",
            "params": { "comparison": ">=", "exclusive": false, "limit": 2 },
            "schemaPath": "#/minimum"
          }]);
        });
      });
    });
  });
});