import React from 'react'
import Markdown from '../Markdown'
import 'highlightjs/styles/github.css'

interface Props {
  opRef: string
  value: any
}

export default (
  {
    opRef,
    value
  }: Props
) =>
  <div className='custom-control-description'>
    <Markdown value={value} opRef={opRef} />
  </div>
