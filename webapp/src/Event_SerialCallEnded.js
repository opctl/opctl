import React from 'react';

export default function EventSerialCallEnded(props) {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      SerialCallEnded
      Id='{props.serialCallEnded.callId}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
