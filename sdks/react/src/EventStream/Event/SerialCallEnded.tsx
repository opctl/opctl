import React from 'react'
import SerialCallEnded from '@opctl/sdk/lib/model/event/serialCallEnded'

interface Props {
  serialCallEnded: SerialCallEnded
  timestamp: Date
}

export default (
  {
    serialCallEnded,
    timestamp
  }: Props
) => {
  return (
    <div style={{ color: 'rgb(96, 253, 255)' }}>
      SerialCallEnded
      Id='{serialCallEnded.callId}'
      Timestamp='{timestamp}'
    </div>
  )
}
