import React from 'react'
import Ansi from 'ansi-to-react'

export default ({
  containerStdOutWrittenTo
}) => {
  return (
    <div>
      <Ansi>
        {window.atob(containerStdOutWrittenTo.data)}
      </Ansi>
    </div>
  )
}
