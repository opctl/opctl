import React from 'react'

interface Props {
  opRef: string
  value: any
}

export default (
  {
    value,
    opRef
  }: Props
) => {
  if (value && typeof value === 'string') {
    value = value.replace(
      /^\/.+$/,
      match => {
        const contentPath = match.slice(0, match.length)
        return `/api/ops/${encodeURIComponent(opRef)}/contents/%2f${encodeURIComponent(contentPath)}`
      })
  } else {
    value = 'opspec-icon.svg'
  }

  return (<img src={value} alt={'icon'} style={{ height: '10vw' }} />)
}
