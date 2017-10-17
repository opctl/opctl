const Ajv = require('ajv');
const ajv = new Ajv();

class Validator {
  constructor() {
    // add formats not included in JSON schema standard
    ajv.addFormat("integer", '^[0-9]+$')
  }

  /**
   * validates value against constraints
   * @param {Number} value
   * @param {Object} [constraints]
   * @return {Array<Error>}
   */
  validate(value,
           constraints) {
    if(!constraints){
      return [];
    }

    ajv.validate(constraints, value);
    return ajv.errors;
  }
}

// export as singleton
module.exports = new Validator();
