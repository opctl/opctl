import React, { CSSProperties } from 'react'
import HasSerialCall from '../HasSerialCall'
import HasParallelCall from '../HasParallelCall'
import HasParallelLoopCall from '../HasParallelLoopCall'
import HasSerialLoopCall from '../HasSerialLoopCall'
import HasCallName from '../HasCallName'
import HasOpCall from '../HasOpCall'
import path from 'path'

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
    if?: Predicate[]
    name?: string
    needs?: string[]
    op?: CallOp
    parallel?: CallParallel
    parallelLoop?: CallParallelLoop
    serial?: CallSerial
    serialLoop?: CallSerialLoop
}

interface Props {
    call: Call
    opRef: string
    style?: CSSProperties
}

export default (
    {
        call,
        opRef
    }: Props
) => {
    let callComponent: any = null

    if (call.op) {
        callComponent = <HasOpCall
            opRef={
                path.join(call.op.ref.startsWith('.') ? path.dirname(opRef) : '', call.op.ref, 'op.yml')
            }
        />
    } else if (call.parallel) {
        callComponent = <HasParallelCall
            call={call}
            opRef={opRef}
        />
    } else if (call.parallelLoop) {
        callComponent = <HasParallelLoopCall
            call={call}
            opRef={opRef}
        />
    } else if (call.serial) {
        callComponent = <HasSerialCall
            call={call}
            opRef={opRef}
        />
    } else if (call.serialLoop) {
        callComponent = <HasSerialLoopCall
            call={call}
            opRef={opRef}
        />
    } else {
        callComponent = null
    }

    return (
        <div
            style={{
                alignItems: 'center',
                display: 'flex',
                flexDirection: 'column'
            }}
        >
            <HasCallName
                call={call}
            />
            {callComponent}
        </div>
    )
}