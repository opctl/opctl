require('../../model/data')
const nodeApiClient = require('../../node/api/client')

/**
 * Coerces a value to a string
 * @param {Value} value
 */
async function toString (value) {
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
      const fileStream = await nodeApiClient.dataGet({ dataRef: value.file })
      return {
        string: await fileStream.text()
      }
    } catch (err) {
      // don't include value in msg; might be secret
      throw new Error(`unable to coerce file to string; error was ${err.message}`)
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

module.exports = toString
