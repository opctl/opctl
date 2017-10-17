const Ajv = require('ajv');
const ajv = new Ajv();
const dockerImageRefParser = require('../../vendored/@codefresh-io/docker-reference/parsers');

class Validator {
  constructor() {
    // add formats not included in JSON schema standard
    ajv.addFormat("docker-image-ref", instance => {
      try {
        dockerImageRefParser.parseQualifiedName(instance)
      } catch (err) {
        return false;
      }
      return true;
    })
  }

  /**
   * validates value against constraints
   * @param {String} value
   * @param {Object} [constraints]
   * @return {Array<Error>}
   */
  validate(value,
           constraints) {
    constraints = Object.assign({type: 'string'}, constraints);

    ajv.validate(constraints, value);
    return ajv.errors || [];
  }
}

// export as singleton
module.exports = new Validator();
