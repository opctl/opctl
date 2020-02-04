import React from 'react'
import Ansi from 'ansi-to-react'

export default ({
  containerStdOutWrittenTo
}) => {
  return (
    <div>
      <Ansi
        linkify
      >
        {window.atob(containerStdOutWrittenTo.data)}
      </Ansi>
    </div>
  )
}
