import _assertStatusSuccessful from './_assertStatusSuccessful'
import constructDataUrl from '../constructDataUrl'

interface FsEntry {
    Mode: string
    Path: string
    Size: number
}

export default async function get(
    dataRef: string
): Promise<FsEntry[]> {
    return fetch(
        constructDataUrl(dataRef)
    )
        .then(_assertStatusSuccessful)
        .then(response => response.json() as Promise<FsEntry[]>)
}