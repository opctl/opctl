'use strict';

var objectUnderTest = require('./validator');

describe('validate', function () {
  describe('constraints undefined', function () {
    test('returns expected result', function () {
      /* arrange/act */
      var actualErrors = objectUnderTest.validate([]);

      /* assert */
      expect(actualErrors).toEqual([]);
    });
  });
  describe('constraints defined', function () {
    describe('minItems constraint', function () {
      describe('value.length > minItems', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate(['dummyItem'], {
            minItems: 1
          });

          /* assert */
          expect(actualErrors).toEqual([]);
        });
      });
      describe('value.length < minItems', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate([], {
            minItems: 1
          });

          /* assert */
          expect(actualErrors).toEqual([{
            "dataPath": "",
            "keyword": "minItems",
            "message": "should NOT have less than 1 items",
            "params": { "limit": 1 },
            "schemaPath": "#/minItems"
          }]);
        });
      });
    });
  });
});