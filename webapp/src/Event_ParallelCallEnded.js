import React from 'react';

export default function EventParallelCallEnded(props) {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      ParallelCallEnded
      Id='{props.parallelCallEnded.callId}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
