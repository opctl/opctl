import React from 'react'
import OpStarted from '@opctl/sdk/lib/model/event/opStarted'

interface Props {
  opStarted: OpStarted
  timestamp: Date
}

export default (
  {
    opStarted,
    timestamp
  }: Props
) => {
  return (
    <div style={{ color: 'rgb(96, 253, 255)' }}>
      OpStarted
      Id='{opStarted.opId}'
      OpRef='{opStarted.opRef}'
      Timestamp='{timestamp}'
    </div>
  )
}
