export default interface ContainerStdOutWrittenTo {
    containerId: string
    data: Int8Array
    imageRef: string
    opRef: string
    rootOpId: string
}