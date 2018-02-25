import React from 'react';

export default ({
                  opStarted,
                  timestamp,
                }) => {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      OpStarted
      Id='{opStarted.opId}'
      PkgRef='{opStarted.pkgRef}'
      Timestamp='{timestamp}'
    </div>
  );
}
