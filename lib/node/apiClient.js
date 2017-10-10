const axios = require('axios');

class ApiClient {
  constructor({baseUrl = 'http://localhost:42224'}){
    this.baseUrl = baseUrl;
  }

  /**
   * Starts an op
   * @param opStartReq
   * @return {Promise.<String>} id of the started op
   */
  async op_start(opStartReq) {
    return axios.post(`${this.baseUrl}/ops/starts`, opStartReq)
      .then(response => (response.data))
  }

  /**
   * Kills an op
   * @param opKillReq
   * @return {Promise.<null>}
   */
  async op_kill(opKillReq) {
    return axios.post(`${this.baseUrl}/ops/kills`, opKillReq)
      .then(null);
  }

  /**
   * Gets pkg content at contentPath
   * @param pkgRef
   * @param contentPath
   * @return {Promise.<String>}
   */
  async pkg_content_get({pkgRef, contentPath}) {
    return axios.get(
      `${this.baseUrl}/pkgs/${encodeURIComponent(pkgRef)}/contents/${encodeURIComponent(contentPath)}`, {
        responseType: 'text'
      })
      .then(response => (response.data));
  }
}

module.exports = ApiClient;
