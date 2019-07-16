/**
 * validates value against constraints
 * @param {boolean} value
 * @return {Array<Error>}
 */
export default function validate(
  value: boolean
) {
  if (typeof value !== 'boolean') {
    return [new Error('boolean required')]
  }

  return []
}