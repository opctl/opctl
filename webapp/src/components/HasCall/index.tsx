import React, { CSSProperties } from 'react'
import CallHasSerial from '../CallHasSerial'
import CallHasParallel from '../CallHasParallel'
import CallHasParallelLoop from '../CallHasParallelLoop'
import CallHasSerialLoop from '../CallHasSerialLoop'
import CallHasName from '../CallHasName'
import CallHasOp from '../CallHasOp'

export interface PullCreds {
    username: string
    password: string
}

export interface CallContainerImage {
    ref: string
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
    inputs?: CallOpInputs
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
    parentOpRef: string
    style?: CSSProperties
}

export default (
    {
        call,
        parentOpRef
    }: Props
) => {
    let callComponent: any = null

    if (call.op) {
        callComponent = <CallHasOp
            callOp={call.op}
            parentOpRef={parentOpRef}
        />
    } else if (call.parallel) {
        callComponent = <CallHasParallel
            callParallel={call.parallel}
            parentOpRef={parentOpRef}
        />
    } else if (call.parallelLoop) {
        callComponent = <CallHasParallelLoop
            callParallelLoop={call.parallelLoop}
            parentOpRef={parentOpRef}
        />
    } else if (call.serial) {
        callComponent = <CallHasSerial
            callSerial={call.serial}
            parentOpRef={parentOpRef}
        />
    } else if (call.serialLoop) {
        callComponent = <CallHasSerialLoop
            callSerialLoop={call.serialLoop}
            parentOpRef={parentOpRef}
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
            <CallHasName
                call={call}
            />
            {callComponent}
        </div>
    )
}