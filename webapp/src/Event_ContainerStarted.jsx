import React from 'react';

export default function EventContainerStarted(props) {
  return (
    <div style={{color: 'rgb(96, 253, 255)'}}>
      ContainerStarted
      Id='{props.containerStarted.containerId}'
      PkgRef='{props.containerStarted.pkgRef}'
      Timestamp='{props.timestamp}'
    </div>
  );
}
