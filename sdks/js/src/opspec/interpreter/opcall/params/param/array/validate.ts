import Ajv from 'ajv'
const ajv = new Ajv()

/**
* validates value against constraints
* @param {Array} value
* @param {Object} [constraints]
* @return {Array<Error>}
*/
export default function validate(
 value: any[],
 constraints?: any
) {
 constraints = Object.assign({ type: 'array' }, constraints)

 ajv.validate(constraints, value)
 return ajv.errors || []
}