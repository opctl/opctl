const keyQualifier = 'content';
class ContentStore {
  /**
   * gets content from store
   * @param key
   * @returns {*}
   */
  get({key}){
    if (global.localStorage) {
      const qualifiedKey = [keyQualifier, key].join('.');
      return JSON.parse(global.localStorage.getItem(qualifiedKey)) || null;
    }
    return null;
  }

  /**
   * sets content in store
   * @param key
   * @param value
   */
  set({key, value}){
    if (global.localStorage) {
      const qualifiedKey = [keyQualifier, key].join('.');
      global.localStorage.setItem(
        qualifiedKey,
        JSON.stringify(value)
      );
    }
  }
}

export default new ContentStore();
