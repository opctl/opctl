import React from 'react'

export default ({
  opErred,
  timestamp
}) => {
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
