const keyQualifier = 'content'
class ContentStore {
  /**
   * gets content from store
   * @param key
   * @returns {*}
   */
  get ({ key }) {
    if (window.localStorage) {
      const qualifiedKey = [keyQualifier, key].join('.')
      const item = window.localStorage.getItem(qualifiedKey)
      if(item) {
        return JSON.parse(item) || null
      }
    }
    return null
  }

  /**
   * sets content in store
   * @param key
   * @param value
   */
  set ({ key, value }) {
    if (window.localStorage) {
      const qualifiedKey = [keyQualifier, key].join('.')
      window.localStorage.setItem(
        qualifiedKey,
        JSON.stringify(value)
      )
    }
  }
}

export default new ContentStore()
