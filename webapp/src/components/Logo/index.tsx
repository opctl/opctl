import React, { CSSProperties } from 'react'

interface Props {
  style?: CSSProperties
}

export default function Logo(
  {
    style
  }: Props
) {
  return (
    <a
      style={style}
      href={`//${window.location.host}`}
    >
      <img
        alt='opctl logo'
        style={{
          height: '2rem'
        }}
        src="/logo.svg"
      />
    </a>
  )
}
