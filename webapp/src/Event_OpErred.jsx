import React from 'react';

export default function EventOpErred(props) {
  return (
    <div style={{color: 'rgb(255, 110, 103)'}}>
      OpErred
      Id='{props.opErred.opId}'
      PkgRef='{props.opErred.pkgRef}'
      Timestamp='{props.timestamp}'
      Msg='{props.opErred.msg}'
    </div>
  );
}
