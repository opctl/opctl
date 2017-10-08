import jsYaml from 'js-yaml';
import opspecNodeApiClient from './opspecNodeApiClient';

class PkgFetcher {
  async fetch(pkgRef) {
    return opspecNodeApiClient.getPkgContent({
      pkgRef,
      contentPath: '/op.yml',
    })
      .then(PkgFetcher._parse)
  }

  static _parse(opDotYml) {
    return jsYaml.safeLoad(opDotYml);
  }
}

// export as singleton
export default new PkgFetcher();
