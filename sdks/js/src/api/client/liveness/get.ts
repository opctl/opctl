import nodeFetch from 'node-fetch'
import _assertStatusSuccessful from '../_assertStatusSuccessful'

/**
 * Gets liveness of node
 *
 * implements https://github.com/opctl/spec/blob/0.1.5/spec/node-api.spec.yml#L60
 */
export default async function get(
    apiBaseUrl: string
): Promise<string> {
    return nodeFetch(
        `${apiBaseUrl}/liveness`
    )
        .then(_assertStatusSuccessful)
        .then(response => (response.text()))
}