import React from 'react';

export default ({
                  containerExited,
                  timestamp,
                }) => {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      ContainerExited
      Id='{containerExited.containerId}'
      OpRef='{containerExited.opRef}'
      ExitCode='{containerExited.exitCode}'
      Timestamp='{timestamp}'
    </div>
  );
}
