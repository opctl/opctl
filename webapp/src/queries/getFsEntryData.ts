import _assertStatusSuccessful from './_assertStatusSuccessful'
import path from 'path'

export default async function get(
    dataRef: string
): Promise<string> {
    return fetch(
        `http://localhost:42224/api/pkgs/${encodeURIComponent(path.dirname(dataRef))}/contents/${encodeURIComponent(path.basename(dataRef))}`
    )
        .then(_assertStatusSuccessful)
        .then(response => response.text())
}