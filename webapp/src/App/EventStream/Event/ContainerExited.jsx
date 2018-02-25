import React from 'react';

export default ({
                  containerExited,
                  timestamp,
                }) => {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      ContainerExited
      Id='{containerExited.containerId}'
      PkgRef='{containerExited.pkgRef}'
      ExitCode='{containerExited.exitCode}'
      Timestamp='{timestamp}'
    </div>
  );
}
