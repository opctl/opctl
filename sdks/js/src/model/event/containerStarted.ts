export default interface ContainerStarted {
    containerId: string
    rootCallId: string
    opRef: string
    name?: string
}
