import jsYaml from 'js-yaml';
import opspecNodeApiClient from './clients/opspecNodeApi';

class PkgFetcher {
  async fetch(pkgRef) {

    return opspecNodeApiClient.pkg_content_get({
      pkgRef,
      contentPath: '/op.yml',
    })
      .then(data => data.text())
      .then(PkgFetcher._parse)
  }

  static _parse(opDotYml) {
    return jsYaml.safeLoad(opDotYml);
  }
}

// export as singleton
export default new PkgFetcher();
