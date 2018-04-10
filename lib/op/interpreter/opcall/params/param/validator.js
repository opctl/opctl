require('../../../../../model/data')
require('../../../../../model/param')
const arrayValidator = require('./array/validator')
const booleanValidator = require('./boolean/validator')
const dirValidator = require('./dir/validator')
const fileValidator = require('./file/validator')
const numberValidator = require('./number/validator')
const objectValidator = require('./object/validator')
const socketValidator = require('./socket/validator')
const stringValidator = require('./string/validator')

class Validator {
  /**
   * Validates value against param
   * note: param defaults aren't considered
   * @param {Value} value
   * @param {Param} param
   * @returns {Array<String>}
   */
  validate (
    value,
    param
  ) {
    if (!param) {
      throw new Error('param required')
    }

    if (param.array) {
      return arrayValidator.validate(
        value.array,
        param.array.constraints
      )
    }

    if (param.boolean) {
      return booleanValidator.validate(
        value.boolean
      )
    }

    if (param.dir) {
      return dirValidator.validate(
        value.dir
      )
    }

    if (param.file) {
      return fileValidator.validate(
        value.file
      )
    }

    if (param.number) {
      return numberValidator.validate(
        value.number,
        param.number.constraints
      )
    }

    if (param.object) {
      return objectValidator.validate(
        value.object,
        param.object.constraints
      )
    }

    if (param.socket) {
      return socketValidator.validate(
        value.socket
      )
    }

    if (param.string) {
      return stringValidator.validate(
        value.string,
        param.string.constraints
      )
    }
  }
}

// export as singleton
module.exports = new Validator()
