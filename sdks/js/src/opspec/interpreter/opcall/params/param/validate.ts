import Value from '../../../../../model/value'
import Param from '../../../../../model/param'
import arrayValidate from './array/validate'
import booleanValidate from './boolean/validate'
import dirValidate from './dir/validate'
import fileValidate from './file/validate'
import numberValidate from './number/validate'
import objectValidate from './object/validate'
import socketValidate from './socket/validate'
import stringValidate from './string/validate'

  /**
   * Validates value against param
   * note: param defaults aren't considered
   * @param {Value} value
   * @param {Param} param
   * @returns {Array<String>}
   */
  export default function validate (
    value: Value,
    param: Param
  ) {
    if (!param) {
      throw new Error('param required')
    }

    if (param.array) {
      return arrayValidate(
        value.array!,
        param.array.constraints
      )
    }

    if (param.boolean) {
      return booleanValidate(
        value.boolean!
      )
    }

    if (param.dir) {
      return dirValidate(
        value.link
      )
    }

    if (param.file) {
      return fileValidate(
        value.link
      )
    }

    if (param.number) {
      return numberValidate(
        value.number!,
        param.number.constraints
      )
    }

    if (param.object) {
      return objectValidate(
        value.object!,
        param.object.constraints
      )
    }

    if (param.socket) {
      return socketValidate(
        value.socket!
      )
    }

    if (param.string) {
      return stringValidate(
        value.string!,
        param.string.constraints
      )
    }
  }
