import getOp from './getOp'

jest.mock('./getFsEntryData')
import getFsEntryData from './getFsEntryData'
const getData = getFsEntryData as jest.MockedFunction<typeof getFsEntryData>

const originalFetch = global.fetch
beforeEach(() => { getData.mockReset() })
afterEach(() => { global.fetch = originalFetch })

test('keeps YAML loading and only describes when op.yml is absent', async () => {
  getData.mockResolvedValue('name: yaml\nrun: {container: {image: {ref: alpine}}}')
  const fetch = jest.fn()
  global.fetch = fetch
  expect((await getOp('/ops/yaml')).name).toBe('yaml')
  expect(fetch).not.toHaveBeenCalled()

  getData.mockRejectedValue(Object.assign(new Error('missing'), { response: { status: 404 } }))
  fetch.mockResolvedValue({ status: 200, json: async () => ({ name: 'go', inputs: { env: { string: {} } } }) })
  const op = await getOp('/ops/go')
  expect(op.name).toBe('go')
  expect(op.run).toBeUndefined()
  expect(fetch).toHaveBeenCalledWith('http://localhost:42224/api/ops/describes', {
    method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ op: { ref: '/ops/go' } })
  })
})

test('keeps authentication errors on the YAML path', async () => {
  const failure = Object.assign(new Error('unauthenticated'), { response: { status: 401 } })
  getData.mockRejectedValue(failure)
  const fetch = jest.fn()
  global.fetch = fetch
  await expect(getOp('/ops/private')).rejects.toBe(failure)
  expect(fetch).not.toHaveBeenCalled()
})
