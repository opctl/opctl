import path from 'path'

export default function constructDataUrl(
    dataRef: string
): string {
    return `http://localhost:42224/api/data/${encodeURIComponent(path.dirname(dataRef))}/contents/${encodeURIComponent(path.basename(dataRef))}`
}
