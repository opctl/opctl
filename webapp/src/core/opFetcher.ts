import jsYaml from 'js-yaml'
import opspecNodeApiClient from './clients/opspecNodeApi'

class OpFetcher {
  async fetch (opRef) {
    return opspecNodeApiClient.data_get({
      dataRef: `${opRef}/op.yml`
    })
      .then(data => data.text())
      .then(OpFetcher._parse)
  }

  static _parse (opDotYml) {
    return jsYaml.safeLoad(opDotYml)
  }
}

// export as singleton
export default new OpFetcher()
