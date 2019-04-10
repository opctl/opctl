import Value from '../value'

export default interface SerialCallEndedEvent {
    callId: string
    outputs: { [key: string]: Value }
    rootOpId: string
}