const fetch = require('node-fetch');

class ApiClient {
  constructor({baseUrl = 'http://localhost:42224'}) {
    this.baseUrl = baseUrl;
  }

  /**
   * Starts an op
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L70
   * @param opStartReq
   * @return {Promise.<String>} id of the started op
   */
  op_start(opStartReq) {
    return fetch(`${this.baseUrl}/ops/starts`, {
      method: 'POST',
      body: JSON.stringify(opStartReq),
    })
      .then(response => (response.text()))
  }

  /**
   * Kills an op
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L139
   * @param opKillReq
   * @return {Promise.<null>}
   */
  op_kill(opKillReq) {
    return fetch(`${this.baseUrl}/ops/kills`, {
      method: 'POST',
      body: JSON.stringify(opKillReq),
    })
      .then(() => null);
  }

  /**
   * Gets pkg content at contentPath
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L242
   * @param pkgRef
   * @param contentPath
   * @return {Promise.<fetch.Response>}
   */
  pkg_content_get({pkgRef, contentPath}) {
    return fetch(
      `${this.baseUrl}/pkgs/${encodeURIComponent(pkgRef)}/contents/${encodeURIComponent(contentPath)}`
    );
  }

  /**
   * Lists pkg contents
   *
   * implements https://github.com/opspec-io/spec/blob/0.1.5/spec/node-api.spec.yml#L178
   * @param pkgRef
   * @return {Promise.<Object[]>}}
   */
  pkg_content_list({pkgRef}) {
    return fetch(
      `${this.baseUrl}/pkgs/${encodeURIComponent(pkgRef)}/contents`
    )
      .then(response => (response.json()));
  }
}

module.exports = ApiClient;
