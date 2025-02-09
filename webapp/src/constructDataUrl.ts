export default function constructDataUrl(
    dataRef: string
): string {
    return `http://127.0.42.224/api/data/${encodeURIComponent(dataRef)}`
}