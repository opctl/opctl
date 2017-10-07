import React from 'react';

export default ({
                  opEnded,
                  timestamp,
                }) => {
  return (
    <div style={{color: opEnded.outcome === 'FAILED' ? 'rgb(255, 110, 103)' : 'rgb(95, 250, 104)'}}>
      OpEnded
      Id='{opEnded.opId}'
      PkgRef='{opEnded.pkgRef}'
      Outcome='{opEnded.outcome}'
      Timestamp='{timestamp}'
    </div>
  );
}
