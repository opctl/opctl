import _assertStatusSuccessful from './_assertStatusSuccessful'
import constructDataUrl from '../constructDataUrl'

export default async function getFsEntryData(
    dataRef: string
): Promise<string> {
    return fetch(
        constructDataUrl(dataRef)
    )
        .then(_assertStatusSuccessful)
        .then(response => response.text())
}