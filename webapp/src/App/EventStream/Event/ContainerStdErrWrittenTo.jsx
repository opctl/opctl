import React from 'react'
import Ansi from 'ansi-to-react'

export default ({
  containerStdErrWrittenTo
}) => {
  return (
    <div>
      <Ansi>
        {window.atob(containerStdErrWrittenTo.data)}
      </Ansi>
    </div>
  )
}
