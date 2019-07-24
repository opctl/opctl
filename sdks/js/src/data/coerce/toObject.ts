import Value from '../../types/value'
import dataGet from '../../api/client/data/get'

/**
 * Coerces a value to an object
 * @param {Value} value
 */
export default async function toArray(
  apiBaseUrl: string,
  value: Value
) {
  if (value.array) {
    // don't include value in msg; might contain secrets
    throw new Error(`unable to coerce array to object; incompatible types`)
  }

  if (value.boolean) {
    throw new Error(`unable to coerce boolean '${value.boolean}' to object; incompatible types`)
  }

  if (value.dir) {
    throw new Error(`unable to coerce dir '${value.dir}' to object; incompatible types`)
  }

  if (value.file) {
    try {
      const fileStream = await dataGet(apiBaseUrl, value.file)
      return parseJsonToObjectValue(await fileStream.text())
    } catch (err) {
      // don't include value in msg; might be secret
      throw new Error(`unable to coerce file to array; error was ${err.message}`)
    }
  }

  if (value.number) {
    // don't include value in msg; might be secret
    throw new Error('unable to coerce number to object; incompatible types')
  }

  if (value.object) {
    return value
  }

  if (value.socket) {
    throw new Error(`unable to coerce socket '${value.socket}' to object; incompatible types`)
  }

  if (value.string) {
    try {
      return parseJsonToObjectValue(value.string)
    } catch (err) {
      // don't include value in msg; might be secret
      throw new Error(`unable to coerce string to object; error was ${err.message}`)
    }
  }

  throw new Error(`unable to coerce ${JSON.stringify(value)} to object`)
}

/**
 * parses JSON to an object value
 * @throws {Error} if possibleJsonArray isn't JSON or isn't an object
 * @returns {Value} the parsed object value
 */
function parseJsonToObjectValue(
  possibleJsonArray: string
) {
  const object = JSON.parse(possibleJsonArray)

  if (Array.isArray(object)) {
    // object === Object([]) returns true; explicitly handle array
    throw new Error('parsed array but expected object')
  }

  if (object === Object(object)) {
    return { object }
  }

  if (object === null) {
    // typeof null returns object; explicitly handle null
    throw new Error('parsed null but expected object')
  }
  throw new Error(`parsed ${typeof object} but expected object`)
}
