import React from 'react'
import ContainerStarted from '@opctl/sdk/lib/model/event/containerStarted'

interface Props {
  containerStarted: ContainerStarted
  timestamp: Date
}

export default (
  {
    containerStarted,
    timestamp
  }: Props
) => {
  return (
    <div style={{ color: 'rgb(96, 253, 255)' }}>
      ContainerStarted
      Id='{containerStarted.containerId}'
      OpRef='{containerStarted.opRef}'
      Timestamp='{timestamp}'
    </div>
  )
}
