const Ajv = require('ajv');
const ajv = new Ajv();

class Validator {
  /**
   * validates value against constraints
   * @param {Array} value
   * @param {Object} [constraints]
   * @return {Array<Error>}
   */
  validate(value,
           constraints) {
    constraints = Object.assign({type: 'array'}, constraints);

    ajv.validate(constraints, value);
    return ajv.errors || [];
  }
}

// export as singleton
module.exports = new Validator();
