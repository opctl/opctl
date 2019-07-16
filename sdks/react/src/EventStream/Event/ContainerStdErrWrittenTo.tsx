import React from 'react'
import Ansi from 'ansi-to-react'
import ContainerStdErrWrittenTo from '@opctl/sdk/lib/model/event/containerStdErrWrittenTo'

interface Props {
  containerStdErrWrittenTo: ContainerStdErrWrittenTo
}

export default (
  {
    containerStdErrWrittenTo
  }: Props
) => {
  return (
    <div>
      <Ansi
        linkify
      >
        {window.atob(containerStdErrWrittenTo.data.toString())}
      </Ansi>
    </div>
  )
}
