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
    describe('format constraint', function () {
      describe('constraint.format is integer', function () {
        describe('value is integer', function () {
          test('returns expected result', function () {
            /* arrange/act */
            var actualErrors = objectUnderTest.validate('1', {
              format: 'integer'
            });

            /* assert */
            expect(actualErrors).toBeNull();
          });
        });
        describe('value isnt integer', function () {
          test('returns expected result', function () {
            /* arrange/act */
            var actualErrors = objectUnderTest.validate('1.1', {
              format: 'integer'
            });

            /* assert */
            expect(actualErrors).toEqual([{
              "dataPath": "",
              "keyword": "format",
              "message": 'should match format "integer"',
              "params": { "format": "integer" },
              "schemaPath": "#/format"
            }]);
          });
        });
      });
    });
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