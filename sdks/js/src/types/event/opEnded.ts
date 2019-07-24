import Value from "../value";

export default interface OpEnded {
    opId: string
    opRef: string
    outcome: 'SUCCEEDED' | 'FAILED' | 'KILLED'
    outputs: { [key: string]: Value }
    rootOpId: string
}