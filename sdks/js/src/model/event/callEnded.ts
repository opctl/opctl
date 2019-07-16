import Value from '../value'

export default interface CallEnded {
    callId: string
    error?: {
        message: string
    }
    outputs: { [key: string]: Value }
    rootOpId: string
}