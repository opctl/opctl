import React from 'react';

export default function EventContainerExited(props) {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      ContainerExited
      Id='{props.containerExited.containerId}'
      PkgRef='{props.containerExited.pkgRef}'
      ExitCode='{props.containerExited.exitCode}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
