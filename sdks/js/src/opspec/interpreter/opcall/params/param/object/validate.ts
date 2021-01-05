import Ajv from 'ajv'
const ajv = new Ajv()

/**
 * validates value against constraints
 * @param {Object} value
 * @param {Object} [constraints]
 * @return {Array<Error>}
 */
export default function validate (
  value?: object | string,
  constraints?: any
) {
  constraints = Object.assign({type: 'object'}, constraints)

  ajv.validate(constraints, value)
  return ajv.errors || []
}
