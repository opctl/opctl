class Validator {
  /**
   * validates value against constraints
   * @param {boolean} value
   * @return {Array<Error>}
   */
  validate (value) {
    if (typeof value !== 'boolean') {
      return [new Error('boolean required')]
    }

    return []
  }
}

// export as singleton
module.exports = new Validator()
