import React, { CSSProperties } from 'react'
import HasContainerCall from '../HasContainerCall'
import HasOpCall from '../HasOpCall'
import HasSerialCall from '../HasSerialCall'
import HasParallelCall from '../HasParallelCall'
import HasParallelLoopCall from '../HasParallelLoopCall'
import HasSerialLoopCall from '../HasSerialLoopCall'

export interface PullCreds {
    username: string
    password: string
}

export interface CallContainerImage {
    ref?: string
    src?: string
    pullCreds?: PullCreds
}

export interface CallContainer {
    image: CallContainerImage
    cmd: string[]
}

export interface CallOpInputs {
    [key: string]: any
}

export interface CallOp {
    ref: string
    inputs: CallOpInputs
}

export type CallParallel = Call[]

export interface LoopVars {
    key?: string
    index?: string
    value?: string
}

export type LoopRangeArray = [string | number]

export type LoopRangeObject = {
    [key: string]: any
}

export type LoopRange = LoopRangeArray | LoopRangeObject | string

export interface CallParallelLoop {
    range: LoopRange
    vars?: LoopVars
    run: Call
}

export interface Predicate {
    eq?: any[]
    ne?: any[]
    exists?: string
    notExists?: string
}

export interface CallSerialLoop {
    range: LoopRange
    vars?: LoopVars
    run: Call
    until?: Predicate[] 
}

export type CallSerial = Call[]

export interface Call {
    container?: CallContainer
    op?: CallOp
    parallel?: CallParallel
    parallelLoop?: CallParallelLoop
    serial?: CallSerial
    serialLoop?: CallSerialLoop
}

interface Props {
    call: Call
    style?: CSSProperties
}

export default (
    {
        call
    }: Props
) => {
    let callComponent: any = null

    if (call.container) {
        callComponent = <HasContainerCall
            call={call}
        />
    } else if (call.op) {
        callComponent = <HasOpCall
            call={call}
        />
    } else if (call.parallel) {
        callComponent = <HasParallelCall
            call={call}
        />
    } else if (call.parallelLoop) {
        callComponent = <HasParallelLoopCall
            call={call}
        />
    } else if (call.serial) {
        callComponent = <HasSerialCall
            call={call}
        />
    } else if (call.serialLoop) {
        callComponent = <HasSerialLoopCall
            call={call}
        />
    } else {
        throw new Error(`unexpected call ${JSON.stringify(call)}`)
    }

    return (
        <div
            style={{
                alignItems: 'center',
                display: 'flex',
                flexDirection: 'column'
            }}
        >
            {callComponent}
        </div>
    )
}