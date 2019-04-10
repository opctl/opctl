  /**
   * validates value against constraints
   * @param {String} value
   * @return {Array<Error>}
   */
  export default function validate (
    value: string
  ) {
    if (!value) {
      return [new Error('dir required')]
    }

    return []
  }