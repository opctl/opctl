import React from 'react'
import Ansi from 'ansi-to-react'

export default ({
  containerStdErrWrittenTo
}) => {
  return (
    <div>
      <Ansi
        linkify
      >
        {window.atob(containerStdErrWrittenTo.data)}
      </Ansi>
    </div>
  )
}
