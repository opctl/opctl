'use strict';

var objectUnderTest = require('./validator');

describe('validate', function () {
  describe('constraints undefined', function () {
    test('returns expected result', function () {
      /* arrange/act */
      var actualErrors = objectUnderTest.validate('dummyValue');

      /* assert */
      expect(actualErrors).toEqual([]);
    });
  });
  describe('constraints defined', function () {
    describe('minLength constraint', function () {
      describe('value > minLength', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate('dummyValue', {
            minLength: 1
          });

          /* assert */
          expect(actualErrors).toBeNull();
        });
      });
      describe('value < minLength', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate('dummyValue', {
            minLength: 100
          });

          /* assert */
          expect(actualErrors).toEqual([{
            "dataPath": "",
            "keyword": "minLength",
            "message": "should NOT be shorter than 100 characters",
            "params": { "limit": 100 },
            "schemaPath": "#/minLength"
          }]);
        });
      });
    });
  });
});