import React from 'react'

export default ({
  containerStarted,
  timestamp
}) => {
  return (
    <div style={{ color: 'rgb(96, 253, 255)' }}>
      ContainerStarted
      Id='{containerStarted.containerId}'
      OpRef='{containerStarted.opRef}'
      Timestamp='{timestamp}'
    </div>
  )
}
