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
    describe('format constraint', function () {
      describe('constraint.format is docker-image-ref', function () {
        describe('value is docker-image-ref', function () {
          test('returns expected result', function () {
            /* arrange/act */
            var actualErrors = objectUnderTest.validate('owner/registry:tag', {
              format: 'docker-image-ref'
            });

            /* assert */
            expect(actualErrors).toEqual([]);
          });
        });
        describe('value isnt docker-image-ref', function () {
          test('returns expected result', function () {
            /* arrange/act */
            var actualErrors = objectUnderTest.validate('#', {
              format: 'docker-image-ref'
            });

            /* assert */
            expect(actualErrors).toEqual([{
              "dataPath": "",
              "keyword": "format",
              "message": 'should match format "docker-image-ref"',
              "params": { "format": "docker-image-ref" },
              "schemaPath": "#/format"
            }]);
          });
        });
      });
    });
    describe('minLength constraint', function () {
      describe('value > minLength', function () {
        test('returns expected result', function () {
          /* arrange/act */
          var actualErrors = objectUnderTest.validate('dummyValue', {
            minLength: 1
          });

          /* assert */
          expect(actualErrors).toEqual([]);
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