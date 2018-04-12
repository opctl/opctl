require('../../model/data')
const nodeApiClient = require('../../node/api/client')

/**
 * Coerces a value to a number
 * @param {Value} value
 */
async function toNumber (value) {
  if (value.array) {
    // don't include value in msg; might contain secrets
    throw new Error(`unable to coerce array to number; incompatible types`)
  }

  if (value.boolean) {
    throw new Error(`unable to coerce boolean '${value.boolean}' to number; incompatible types`)
  }

  if (value.dir) {
    throw new Error(`unable to coerce dir '${value.dir}' to number; incompatible types`)
  }

  if (value.file) {
    try {
      const fileStream = await nodeApiClient.dataGet({ dataRef: value.file })
      return parseJsonToNumberValue(await fileStream.text())
    } catch (err) {
      // don't include value in msg; might contain secrets
      throw new Error(`unable to coerce file to number; error was ${err.message}`)
    }
  }

  if (value.number) {
    return value
  }

  if (value.object) {
    // don't include value in msg; might contain secrets
    throw new Error(`unable to coerce object to number; incompatible types`)
  }

  if (value.socket) {
    throw new Error(`unable to coerce socket '${value.socket}' to number; incompatible types`)
  }

  if (value.string) {
    try {
      return parseJsonToNumberValue(value.string)
    } catch (err) {
      // don't include value in msg; might be secret
      throw new Error(`unable to coerce string to number; error was ${err.message}`)
    }
  }

  throw new Error(`unable to coerce ${JSON.stringify(value)} to number`)
}

/**
 * parses JSON to a number value
 * @throws {Error} if possibleJsonNumber isn't JSON or isn't a number
 * @returns {Value} the parsed array value
 */
function parseJsonToNumberValue (possibleJsonNumber) {
  const number = JSON.parse(possibleJsonNumber)
  if (typeof number === 'number') {
    return { number }
  }
  throw new Error(`parsed ${typeof number} but expected number`)
}

module.exports = toNumber
