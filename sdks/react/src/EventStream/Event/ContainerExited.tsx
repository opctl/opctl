import React from 'react'
import ContainerExited from '@opctl/sdk/src/model/event/containerExited'

interface Props {
  containerExited: ContainerExited
  timestamp: Date
}

export default (
  {
    containerExited,
    timestamp
  }: Props
) => {
  return (
    <div style={{ color: 'rgb(96, 253, 255)' }}>
      ContainerExited
      Id='{containerExited.containerId}'
      OpRef='{containerExited.opRef}'
      ExitCode='{containerExited.exitCode}'
      Timestamp='{timestamp}'
    </div>
  )
}
