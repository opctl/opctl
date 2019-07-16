export default interface ContainerStdErrWrittenToEvent {
    containerId: string
    data: Int8Array
    imageRef: string
    opRef: string
    rootOpId: string
}