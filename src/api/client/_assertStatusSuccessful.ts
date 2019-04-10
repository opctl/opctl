import { Response } from 'node-fetch'

/**
  * Asserts response.status is in the range of successful status codes
  * @param response
  * @private
  */
export default async function _assertStatusSuccessful(
    response: Response
): Promise<Response> {
    if (response.status >= 200 && response.status < 300) {
        return Promise.resolve(response)
    } else {
        return response.text().then((errorMsg: string) => {
            const error = new Error(errorMsg) as any
            error.response = response
            throw error
        })
    }
}