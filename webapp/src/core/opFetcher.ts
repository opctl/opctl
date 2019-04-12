import jsYaml from 'js-yaml'
import getApiBaseUrl from './getApiBaseUrl'
import { dataGet } from '@opctl/sdk/lib/api/client'

class OpFetcher {
  async fetch(opRef) {
    return dataGet(
      getApiBaseUrl(),
      `${opRef}/op.yml`
    )
      .then(data => data.text())
      .then(OpFetcher._parse)
  }

  static _parse(opDotYml) {
    return jsYaml.safeLoad(opDotYml)
  }
}

// export as singleton
export default new OpFetcher()
