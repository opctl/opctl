import React from 'react'

export default ({
  serialCallEnded,
  timestamp
}) => {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      SerialCallEnded
      Id='{serialCallEnded.callId}'
      Timestamp='{timestamp}'
    </div>
  )
}
