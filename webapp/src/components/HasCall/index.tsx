import React, { CSSProperties } from 'react'
import CallHasSerial from '../CallHasSerial'
import CallHasParallel from '../CallHasParallel'
import CallHasParallelLoop from '../CallHasParallelLoop'
import CallHasSerialLoop from '../CallHasSerialLoop'
import CallHasOp from '../CallHasOp'
import CallHasSummary from '../CallHasSummary'
import brandColors from '../../brandColors'

export interface PullCreds {
  username: string
  password: string
}

export interface ContainerCallImage {
  ref: string
  pullCreds?: PullCreds
}

export interface ContainerCall {
  image: ContainerCallImage
  cmd: string[]
}

export interface OpCallInputs {
  [key: string]: any
}

export interface OpCall {
  ref: string
  inputs?: OpCallInputs
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

export interface ParallelLoopCall {
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

export interface SerialLoopCall {
  range: LoopRange
  vars?: LoopVars
  run: Call
  until?: Predicate[]
}

export type CallSerial = Call[]

export interface Call {
  container?: ContainerCall
  if?: Predicate[]
  name?: string
  needs?: string[]
  op?: OpCall
  parallel?: CallParallel
  parallelLoop?: ParallelLoopCall
  serial?: CallSerial
  serialLoop?: SerialLoopCall
}

interface Props {
  call: Call
  parentOpRef: string
  style?: CSSProperties
}

export default function HasCall(
  {
    call,
    parentOpRef
  }: Props
) {
  let callComponent: any = null

  if (call.op) {
    callComponent = <CallHasOp
      opCall={call.op}
      parentOpRef={parentOpRef}
    />
  } else if (call.parallel) {
    callComponent = <CallHasParallel
      callParallel={call.parallel}
      parentOpRef={parentOpRef}
    />
  } else if (call.parallelLoop) {
    callComponent = <CallHasParallelLoop
      parallelLoopCall={call.parallelLoop}
      parentOpRef={parentOpRef}
    />
  } else if (call.serial) {
    callComponent = <CallHasSerial
      callSerial={call.serial}
      parentOpRef={parentOpRef}
    />
  } else if (call.serialLoop) {
    callComponent = <CallHasSerialLoop
      serialLoopCall={call.serialLoop}
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
        justifyContent: 'center',
        flexDirection: 'column',
        ...!(call.container || call.serial || call.parallel)
          ? {
            border: `solid .1rem ${brandColors.lightGray}`
          }
          : null,
        marginLeft: '1rem',
        marginRight: '1rem'
      }}
    >
      <CallHasSummary
        call={call}
        parentOpRef={parentOpRef}
      />
      {callComponent}
    </div>
  )
}