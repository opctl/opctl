const Ajv = require('ajv');
const ajv = new Ajv();

class Validator {
  validate(value,
           constraints) {
    ajv.validate(constraints, value);
    return ajv.errors;
  }
}

// export as singleton
module.exports = new Validator();
