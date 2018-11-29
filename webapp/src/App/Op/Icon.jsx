import React from 'react'

export default ({ value, opRef }) => {
  if (value) {
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
