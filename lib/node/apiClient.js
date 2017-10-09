const axios = require('axios');

// @TODO: don't assume local node
const baseUrl = 'http://localhost:42224';

class ApiClient {
  /**
   * Starts an op
   * @param startOpReq
   * @return {Promise.<String>} id of the started op
   */
  async op_start(startOpReq) {
    return axios.post(`${baseUrl}/ops/starts`, {
      data: startOpReq,
    })
      .then(response => (response.data))
  }

  /**
   * Gets the pkg content at contentPath
   * @param pkgRef
   * @param contentPath
   * @return {Promise.<String>}
   */
  async pkg_content_get({pkgRef, contentPath}) {
    return axios.get(
      `${baseUrl}/pkgs/${encodeURIComponent(pkgRef)}/contents/${encodeURIComponent(contentPath)}`, {
        responseType: 'text'
      })
      .then(response => (response.data));
  }
}

// export as singleton
module.exports = new ApiClient();
