import path from 'path-browserify'
import jsYaml from 'js-yaml'
import getFsEntryData from './getFsEntryData'
import assertStatusSuccessful from './_assertStatusSuccessful'

export default async function getOp(opRef: string): Promise<any> {
  try {
    return jsYaml.safeLoad(await getFsEntryData(path.join(opRef, 'op.yml')))
  } catch (error) {
    if ((error as any).response?.status !== 404) {
      throw error
    }
    const response = await fetch('http://localhost:42224/api/ops/describes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ op: { ref: opRef } })
    })
    await assertStatusSuccessful(response)
    return response.json()
  }
}
