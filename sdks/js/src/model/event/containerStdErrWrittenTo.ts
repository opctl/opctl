export default interface ContainerStdErrWrittenTo {
    containerId: string
    data: Int8Array
    imageRef: string
    opRef: string
    rootOpId: string
}