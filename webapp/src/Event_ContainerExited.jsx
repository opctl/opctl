import React from 'react';

export default function EventContainerExited(props) {
  return (
    <div>
      ContainerExited
      Id='{props.containerExited.containerId}'
      PkgRef='{props.containerExited.pkgRef}'
      ExitCode='{props.containerExited.exitCode}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
