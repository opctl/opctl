import React from 'react'
import OpErred from '@opctl/sdk/src/model/event/opErred'

interface Props {
  opErred: OpErred
  timestamp: Date
}

export default (
  {
    opErred,
    timestamp
  }: Props
) => {
  return (
    <div style={{ color: 'rgb(255, 110, 103)' }}>
      OpErred
      Id='{opErred.opId}'
      OpRef='{opErred.opRef}'
      Timestamp='{timestamp}'
      Msg='{opErred.msg}'
    </div>
  )
}
