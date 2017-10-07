import React from 'react';

export default ({
                  parallelCallEnded,
                  timestamp,
                }) => {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      ParallelCallEnded
      Id='{parallelCallEnded.callId}'
      Timestamp='{timestamp}'
    </div>
  );
}
