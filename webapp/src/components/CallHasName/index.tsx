import React, { Fragment } from 'react'
import { Call } from '../HasCall'

interface Props {
  call: Call
}

function getName(call: Call): string | null {
  if (call.name) {
    return call.name
  } else if (call.container) {
    return call.container.image?.ref || 'Container'
  } else if (call.op) {
    return call.op.ref || 'Op'
  } else if (call.serialLoop) {
    return 'Serial Loop'
  } else if (call.parallelLoop) {
    return 'Parallel Loop'
  }
  return null
}

export default function CallHasName(
  {
    call
  }: Props
) {
  const name = getName(call)
  if (!name) {
    return null
  }

  return (
    <Fragment>
      {name}
    </Fragment>
  )
}