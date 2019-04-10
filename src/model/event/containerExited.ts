import Value from '../value'

export default interface ContainerExited {
    containerId: string
    exitCode: number
    imageRef: string
    opRef: string
    outputs: { [key: string]: Value }
    rootOpId: string
}