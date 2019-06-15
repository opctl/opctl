import Ajv from 'ajv'
import dockerImageRefParser from '@codefresh-io/docker-reference/lib/parsers'

const ajv = new Ajv()

// add formats not included in JSON schema standard
ajv.addFormat('docker-image-ref', instance => {
  try {
    dockerImageRefParser.parseQualifiedName(instance)
  } catch (err) {
    return false
  }
  return true
})

/**
 * validates value against constraints
 * @param {String} value
 * @param {Object} [constraints]
 * @return {Array<Error>}
 */
export default function validate (
  value: string,
  constraints?: any
  ) {
  constraints = Object.assign({type: 'string'}, constraints)

  ajv.validate(constraints, value)
  return ajv.errors || []
}