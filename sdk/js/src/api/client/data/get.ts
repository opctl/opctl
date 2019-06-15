import _assertStatusSuccessful from '../_assertStatusSuccessful'
import nodeFetch, { Response } from 'node-fetch'

/**
 * Gets data
 *
 * implements https://github.com/opctl/spec/blob/0.1.5/spec/node-api.spec.yml#L242
 * @param dataRef
 */
export default async function get(
    apiBaseUrl: string,
    dataRef: string
): Promise<Response> {
    return nodeFetch(
        `${apiBaseUrl}/data/${encodeURIComponent(dataRef)}`
    )
        .then(_assertStatusSuccessful)
}