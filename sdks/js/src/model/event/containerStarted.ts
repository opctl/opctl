export default interface ContainerStarted {
    containerId: string
    imageRef: string
    opRef: string
    rootOpId: string
}