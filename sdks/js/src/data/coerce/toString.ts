import Value from '../../model/value'
import dataGet from '../../api/client/data/get'

/**
 * Coerces a value to a string
 * @param {Value} value
 */
export default async function toString(
  apiBaseUrl: string,
  value: Value
) {
  if (value.array) {
    return {
      string: JSON.stringify(value.array)
    }
  }

  if (typeof value.boolean === 'boolean') {
    return {
      string: JSON.stringify(value.boolean)
    }
  }

  if (value.file) {
    try {
      const fileStream = await dataGet(apiBaseUrl, value.file)
      return {
        string: await fileStream.text()
      }
    } catch (err) {
      // don't include value in msg; might be secret
      throw new Error(`unable to coerce file to string; error was ${(err as Error).message}`)
    }
  }

  if (value.number) {
    return {
      string: JSON.stringify(value.number)
    }
  }

  if (value.object) {
    return {
      string: JSON.stringify(value.object)
    }
  }

  if (value.string) {
    return value
  }

  throw new Error(`unable to coerce ${JSON.stringify(value)} to string`)
}