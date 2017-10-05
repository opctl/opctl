import React from 'react';

export default function EventOpEnded(props) {
  return (
    <div style={{color: props.opEnded.outcome === 'FAILED' ? 'rgb(255, 110, 103)' : 'rgb(95, 250, 104)'}}>
      OpEnded
      Id='{props.opEnded.opId}'
      PkgRef='{props.opEnded.pkgRef}'
      Outcome='{props.opEnded.outcome}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
