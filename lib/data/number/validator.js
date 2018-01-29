const Ajv = require('ajv');
const ajv = new Ajv();

class Validator {
  constructor() {
    // add formats not included in JSON schema standard
    ajv.addFormat(
      "integer",
      {
        type: 'number',
        validate: value => /^[0-9]+$/.test(value),
      });
  }

  /**
   * validates value against constraints
   * @param {Number} value
   * @param {Object} [constraints]
   * @return {Array<Error>}
   */
  validate(value,
           constraints) {
    constraints = Object.assign({type: 'number'}, constraints);

    ajv.validate(constraints, value);
    return ajv.errors || [];
  }
}

// export as singleton
module.exports = new Validator();
