require('../../model/data')
const nodeApiClient = require('../../node/api/client')

/**
 * Coerces a value to an array
 * @param {Value} value
 */
async function toArray (value) {
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
      const fileStream = await nodeApiClient.dataGet({ dataRef: value.file })
      return parseJsonToArrayValue(await fileStream.text())
    } catch (err) {
      throw new Error(`unable to coerce file to array; error was ${err.message}`)
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
      throw new Error(`unable to coerce string to array; error was ${err.message}`)
    }
  }

  throw new Error(`unable to coerce ${JSON.stringify(value)} to array`)
}

/**
 * parses JSON to an array value
 * @throws {Error} if possibleJsonArray isn't JSON or isn't an array
 * @returns {Value} the parsed array value
 */
function parseJsonToArrayValue (possibleJsonArray) {
  const array = JSON.parse(possibleJsonArray)
  if (Array.isArray(array)) {
    return { array }
  }
  throw new Error(`parsed ${typeof array} but expected array`)
}

module.exports = toArray
