/**
 * validates value against constraints
 */
export default function validate(
  value?: string | null
) {
  if (!value) {
    return [new Error('file required')]
  }

  return []
}
