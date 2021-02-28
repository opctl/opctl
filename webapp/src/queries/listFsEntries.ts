import _assertStatusSuccessful from './_assertStatusSuccessful'

interface FsEntry {
    Mode: string
    Path: string
    Size: number
}

export default async function get(
    dataRef: string
): Promise<FsEntry[]> { 
    return fetch(
        `http://localhost:42224/api/data/${encodeURIComponent(dataRef)}/contents`
    )
        .then(_assertStatusSuccessful)
        .then(response => response.json() as Promise<FsEntry[]>)
}
