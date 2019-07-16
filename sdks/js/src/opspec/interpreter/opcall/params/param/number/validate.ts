import Ajv from 'ajv'
const ajv = new Ajv()

// add formats not included in JSON schema standard
ajv.addFormat(
  'integer',
  {
    type: 'number',
    validate: (value:string) => /^[0-9]+$/.test(value)
  } as any
  )

/**
 * validates value against constraints
 * @param {Number} value
 * @param {Object} [constraints]
 * @return {Array<Error>}
 */
export default function validate(
  value: number,
  constraints?: any
) {
  constraints = Object.assign({ type: 'number' }, constraints)

  ajv.validate(constraints,
    value)
  return ajv.errors || []
}