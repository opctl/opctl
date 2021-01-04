import React from 'react'
import HasCall, { ParallelLoopCall } from '../HasCall'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'


interface Props {
  parallelLoopCall: ParallelLoopCall
  parentOpRef: string
}

export default function CallHasParallelLoop(
  {
    parallelLoopCall,
    parentOpRef
  }: Props
) {
  return (
    <div
      style={{
        display: 'flex',
        alignItems: 'center',
        flexDirection: 'column'
      }}
    >
      <div
        style={{
          backgroundColor: brandColors.lightGray,
          minHeight: '2.5rem',
          height: '100%',
          width: '.1rem'
        }}
      ></div>
      <AddCallPopper>
        <PlusIcon
          style={{
            backgroundColor: brandColors.white,
            cursor: 'pointer',
            fill: brandColors.active,
            display: 'block'
          }}
        />
      </AddCallPopper>
      <div
        style={{
          backgroundColor: brandColors.lightGray,
          minHeight: '2.5rem',
          height: '100%',
          width: '.1rem'
        }}
      ></div>
      <HasCall
        call={parallelLoopCall!.run}
        parentOpRef={parentOpRef}
      />
      <div
        style={{
          backgroundColor: brandColors.lightGray,
          minHeight: '2.5rem',
          height: '100%',
          width: '.1rem'
        }}
      ></div>
      <AddCallPopper>
        <PlusIcon
          style={{
            backgroundColor: brandColors.white,
            cursor: 'pointer',
            fill: brandColors.active,
            display: 'block'
          }}
        />
      </AddCallPopper>
      <div
        style={{
          backgroundColor: brandColors.lightGray,
          minHeight: '2.5rem',
          height: '100%',
          width: '.1rem'
        }}
      ></div>
    </div>
  )
}