import Value from '../../../types/value'
import isStringTruthy from './isStringTruthy'
import dataGet from '../../../api/client/data/get'

/**
 * Coerces a value to a boolean
 * @param {Value} value
 */
export default async function toBoolean(
  apiBaseUrl: string,
  value: Value
) {
  if (value.array) {
    return {
      boolean: value.array.length > 0
    }
  }

  if (typeof value.boolean === 'boolean') {
    return value
  }

  if (value.dir) {
    try {
      const dirStream = await dataGet(apiBaseUrl, value.dir)
      const dir = await dirStream.json()
      return { boolean: dir.length > 0 }
    } catch (err) {
      throw new Error(`unable to coerce dir to boolean; error was ${err.message}`)
    }
  }

  if (value.file) {
    try {
      const fileStream = await dataGet(apiBaseUrl, value.file)
      const file = await fileStream.text()
      return { boolean: isStringTruthy(file) }
    } catch (err) {
      throw new Error(`unable to coerce file to boolean; error was ${err.message}`)
    }
  }

  if (value.number) {
    return {
      boolean: value.number !== 0
    }
  }

  if (value.object) {
    return {
      boolean: Object.entries(value.object).length !== 0
    }
  }

  if (value.socket) {
    throw new Error(`unable to coerce socket '${value.socket}' to boolean; incompatible types`)
  }

  if (value.string) {
    return { boolean: isStringTruthy(value.string) }
  }

  throw new Error(`unable to coerce ${JSON.stringify(value)} to boolean`)
}
