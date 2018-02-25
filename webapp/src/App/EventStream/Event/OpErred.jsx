import React from 'react';

export default ({
                  opErred,
                  timestamp,
                }) => {
  return (
    <div style={{color: 'rgb(255, 110, 103)'}}>
      OpErred
      Id='{opErred.opId}'
      PkgRef='{opErred.pkgRef}'
      Timestamp='{timestamp}'
      Msg='{opErred.msg}'
    </div>
  );
}
