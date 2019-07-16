export default interface ContainerStdOutWrittenToEvent {
    containerId: string
    data: Int8Array
    imageRef: string
    opRef: string
    rootOpId: string
}