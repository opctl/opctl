import Value from '../../model/value'
import dataGet from '../../api/client/data/get'

/**
 * Coerces a value to an array
 * @param {Value} value
 */
export default async function toArray(
  apiBaseUrl: string,
  value: Value
) {
  if (value.array) {
    return value
  }

  if (value.boolean) {
    throw new Error(`unable to coerce boolean '${value.boolean}' to array; incompatible types`)
  }

  if (value.dir) {
    throw new Error(`unable to coerce dir '${value.dir}' to array; incompatible types`)
  }

  if (value.file) {
    try {
      const fileStream = await dataGet(apiBaseUrl, value.file)
      return parseJsonToArrayValue(await fileStream.text())
    } catch (err) {
      throw new Error(`unable to coerce file to array; error was ${(err as Error).message}`)
    }
  }

  if (value.number) {
    // don't include value in msg; might be secret
    throw new Error(`unable to coerce number to array; incompatible types`)
  }

  if (value.object) {
    // don't include value in msg; might contain secrets
    throw new Error(`unable to coerce object to array; incompatible types`)
  }

  if (value.socket) {
    throw new Error(`unable to coerce socket '${value.socket}' to array; incompatible types`)
  }

  if (value.string) {
    try {
      return parseJsonToArrayValue(value.string)
    } catch (err) {
      // don't include value in msg; might be secret
      throw new Error(`unable to coerce string to array; error was ${(err as Error).message}`)
    }
  }

  throw new Error(`unable to coerce ${JSON.stringify(value)} to array`)
}

/**
 * parses JSON to an array value
 * @throws {Error} if possibleJsonArray isn't JSON or isn't an array
 * @returns {Value} the parsed array value
 */
function parseJsonToArrayValue(
  possibleJsonArray: string
) {
  const array = JSON.parse(possibleJsonArray)
  if (Array.isArray(array)) {
    return { array }
  }
  throw new Error(`parsed ${typeof array} but expected array`)
}

module.exports = toArray
