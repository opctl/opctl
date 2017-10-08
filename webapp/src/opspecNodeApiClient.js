// @TODO: don't assume local node
const baseUrl = 'http://localhost:42224';

class OpSpecNodeApiClient {
  /**
   * Checks for http status errors; required per: https://github.com/github/fetch/tree/v2.0.3#handling-http-error-statuses
   * @param response
   * @return {*}
   * @private
   */
  static async _checkStatus(response) {
    if (response.status >= 200 && response.status < 300) {
      return response
    } else {
      return response.text().then(msg => {
        return Promise.reject(
          new Error(`encountered error from ${baseUrl}; code: ${response.status}, msg: ${msg}`)
        );
      });
    }
  }

  /**
   * Starts an op
   * @param startOpReq
   * @return {Promise.<String>} id of the started op
   */
  async startOp(startOpReq) {
    return fetch(`${baseUrl}/ops/starts`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(startOpReq),
    })
      .then(OpSpecNodeApiClient._checkStatus)
      .then(response => (response.text()))
  }

  /**
   * Gets the pkg content at contentPath
   * @param pkgRef
   * @param contentPath
   * @return {Promise.<String>}
   */
  async getPkgContent({pkgRef, contentPath}) {
    return fetch(`${baseUrl}/pkgs/${encodeURIComponent(pkgRef)}/contents/${encodeURIComponent(contentPath)}`)
      .then(OpSpecNodeApiClient._checkStatus)
      .then(response => (response.text()));
  }
}

// export as singleton
export default new OpSpecNodeApiClient();
