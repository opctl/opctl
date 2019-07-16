import nodeFetch from 'node-fetch'
import _assertStatusSuccessful from '../../_assertStatusSuccessful'
import Value from '../../../../model/value'

/**
 * Starts an op
 *
 * implements https://github.com/opctl/spec/blob/0.1.5/spec/node-api.spec.yml#L70
 * @return {Promise.<String>} id of the started op
 */
export default async function post(
    apiBaseUrl: string,
    args: { [key: string]: Value },
    op: {
        ref: string
        pullCreds?: {
            username: string
            password: string
        } | null | undefined
    }
): Promise<string> {
    return nodeFetch(`${apiBaseUrl}/ops/starts`, {
        method: 'POST',
        body: JSON.stringify({
            args,
            op
        })
    })
        .then(_assertStatusSuccessful)
        .then(response => (response.text()))
}