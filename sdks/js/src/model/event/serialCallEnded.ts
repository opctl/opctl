import Value from '../value'

export default interface SerialCallEnded {
    callId: string
    outputs: { [key: string]: Value }
    rootOpId: string
}