import React from 'react'
import { Window as WindowContext } from '../WindowContext'
import OpWindow from '../OpWindow'
import CodeWindow from '../CodeWindow'
import path from 'path-browserify'

interface Props {
  window: WindowContext
}

export default function Window(
  {
    window
  }: Props
) {
  if (
    path.basename(window.fsEntry.path) === 'op.yml'
  ) {
    return (
      <OpWindow
        window={window}
      />
    )
  }

  return (
    <CodeWindow
      window={window}
    />
  )
}
