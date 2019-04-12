import Value from '../value'

export default interface CallEnded {
    callId: string
    outputs: { [key: string]: Value }
    rootOpId: string
}