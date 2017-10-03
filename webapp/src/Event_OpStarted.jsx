import React from 'react';

export default function EventOpStarted(props) {
  return (
    <div>
      OpStarted
      Id='{props.opStarted.opId}'
      PkgRef='{props.opStarted.pkgRef}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
