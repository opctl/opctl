import React from 'react'
import { Call } from '../HasCall'
import CallHasIcon from '../CallHasIcon'
import CallHasName from '../CallHasName'

interface Props {
  call: Call
  parentOpRef?: string
}

export default function CallHasSummary(
  {
    call,
    parentOpRef
  }: Props
) {
  if (
    !call.name
    && (
      call.serial
      || call.parallel
    )
  ) {
    //these call types should only render a description if they have a name
    return null
  }

  return (
    <div
      style={{
        boxShadow: '0 .1rem .5rem rgba(27,31,35,.15), 0 0 .1rem rgba(106,115,125,.35)',
        borderRadius: '.5rem',
        minWidth: '5rem',
        minHeight: '3rem',
        padding: '0 .5rem',
        margin: '0 .5rem',
        display: 'flex',
        width: 'fit-content',
        alignItems: 'center',
        justifyContent: 'center',
        backgroundColor: 'white',
        zIndex: 2,
        marginTop: '-1.5rem'
      }}
    >
      <CallHasIcon
        call={call}
        parentOpRef={parentOpRef}
      />
      <CallHasName
        call={call}
      />
    </div>
  )
}