import React from 'react';

export default function EventOpEnded(props) {
  return (
    <div>
      OpEnded
      Id='{props.opEnded.opId}'
      PkgRef='{props.opEnded.pkgRef}'
      Outcome='{props.opEnded.outcome}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
