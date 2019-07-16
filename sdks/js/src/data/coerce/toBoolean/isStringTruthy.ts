/**
 * ensures value isn't:
 * - all "0"'s
 * - ""
 * - "FALSE" (case insensitive)
 * - "F" (case insensitive)
 * @param {string} value
 */
export default function isStringTruthy(
  value: string
) {
  const normalizedValue = value ? value.toUpperCase().replace(/0/g, '') : ''
  switch (normalizedValue) {
    case '':
      return false
    case 'FALSE':
      return false
    case 'F':
      return false
    default:
      return true
  }
}
