'use strict';

var objectUnderTest = require('./validator');

describe('validate', function () {
  describe('constraints undefined', function () {
    test('returns expected result', function () {
      /* arrange/act */
      var actualErrors = objectUnderTest.validate({});

      /* assert */
      expect(actualErrors).toEqual([]);
    });
  });
  describe('constraints defined', function () {
    describe('minProperties constraint', function () {
      describe('Object.keys(value).length > minProperties', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate({ 'dummyProperty': 'dummyValue' }, {
            minProperties: 1
          });

          /* assert */
          expect(actualErrors).toBeNull();
        });
      });
      describe('Object.keys(value).length < minProperties', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate({}, {
            minProperties: 1
          });

          /* assert */
          expect(actualErrors).toEqual([{
            "dataPath": "",
            "keyword": "minProperties",
            "message": "should NOT have less than 1 properties",
            "params": { "limit": 1 },
            "schemaPath": "#/minProperties"
          }]);
        });
      });
    });
  });
});