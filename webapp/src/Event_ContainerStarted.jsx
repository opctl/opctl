import React from 'react';

export default function EventContainerStarted(props) {
  return (
    <div>
      ContainerStarted
      Id='{props.containerStarted.containerId}'
      PkgRef='{props.containerStarted.pkgRef}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
