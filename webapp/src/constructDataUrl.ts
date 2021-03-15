export default function constructDataUrl(
    dataRef: string
): string {
    return `http://localhost:42224/api/data/${encodeURIComponent(dataRef)}`
}