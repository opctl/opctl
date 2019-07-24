import React from 'react'
import Ansi from 'ansi-to-react'
import ContainerStdOutWrittenTo from '@opctl/sdk/src/types/event/containerStdOutWrittenTo'

interface Props {
  containerStdOutWrittenTo: ContainerStdOutWrittenTo
}

export default (
  {
    containerStdOutWrittenTo
  }: Props
) => {
  return (
    <div>
      <Ansi
        linkify
      >
        {window.atob(containerStdOutWrittenTo.data.toString())}
      </Ansi>
    </div>
  )
}
