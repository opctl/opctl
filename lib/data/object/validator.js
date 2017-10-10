const Ajv = require('ajv');
const ajv = new Ajv();

class Validator {
  /**
   * validates value against constraints
   * @param {Object} value
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
