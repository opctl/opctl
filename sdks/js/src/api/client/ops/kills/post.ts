import nodeFetch from 'node-fetch'
import _assertStatusSuccessful from '../../_assertStatusSuccessful'

/**
 * Kills an op
 *
 * implements https://github.com/opctl/spec/blob/0.1.5/spec/node-api.spec.yml#L139
 * @param {Object} opKillReq
 * @param {Object} opKillReq.opId
 */
export default async function post(
    apiBaseUrl: string,
    opId: string
): Promise<null> {
    return nodeFetch(`${apiBaseUrl}/ops/kills`, {
        method: 'POST',
        body: JSON.stringify({
            opId
        })
    })
        .then(_assertStatusSuccessful)
        .then(() => null)
}