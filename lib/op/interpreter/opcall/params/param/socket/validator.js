class Validator {
  /**
   * validates value against constraints
   * @param {String} value
   * @return {Array<Error>}
   */
  validate (value) {
    if (!value) {
      return [new Error('socket required')]
    }

    return []
  }
}

// export as singleton
module.exports = new Validator()
