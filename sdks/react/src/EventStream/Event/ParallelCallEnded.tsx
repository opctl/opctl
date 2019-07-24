import React from 'react'
import ParallelCallEnded from '@opctl/sdk/src/types/event/parallelCallEnded'

interface Props {
  parallelCallEnded: ParallelCallEnded
  timestamp: Date
}

export default (
  {
    parallelCallEnded,
    timestamp
  }: Props
) => {
  return (
    <div style={{ color: 'rgb(96, 253, 255)' }}>
      ParallelCallEnded
      Id='{parallelCallEnded.callId}'
      Timestamp='{timestamp}'
    </div>
  )
}
