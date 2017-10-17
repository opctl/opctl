'use strict';

var _createClass = function () { function defineProperties(target, props) { for (var i = 0; i < props.length; i++) { var descriptor = props[i]; descriptor.enumerable = descriptor.enumerable || false; descriptor.configurable = true; if ("value" in descriptor) descriptor.writable = true; Object.defineProperty(target, descriptor.key, descriptor); } } return function (Constructor, protoProps, staticProps) { if (protoProps) defineProperties(Constructor.prototype, protoProps); if (staticProps) defineProperties(Constructor, staticProps); return Constructor; }; }();

function _classCallCheck(instance, Constructor) { if (!(instance instanceof Constructor)) { throw new TypeError("Cannot call a class as a function"); } }

var Ajv = require('ajv');
var ajv = new Ajv();
var dockerImageRefParser = require('../../vendored/@codefresh-io/docker-reference/parsers');

var Validator = function () {
  function Validator() {
    _classCallCheck(this, Validator);

    // add formats not included in JSON schema standard
    ajv.addFormat("docker-image-ref", function (instance) {
      try {
        dockerImageRefParser.parseQualifiedName(instance);
      } catch (err) {
        return false;
      }
      return true;
    });
  }

  /**
   * validates value against constraints
   * @param {String} value
   * @param {Object} [constraints]
   * @return {Array<Error>}
   */


  _createClass(Validator, [{
    key: 'validate',
    value: function validate(value, constraints) {
      constraints = Object.assign({ type: 'string' }, constraints);

      ajv.validate(constraints, value);
      return ajv.errors || [];
    }
  }]);

  return Validator;
}();

// export as singleton


module.exports = new Validator();