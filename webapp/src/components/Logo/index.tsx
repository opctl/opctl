import React, { CSSProperties } from 'react'

interface Props {
  style?: CSSProperties
}

export default (
  {
    style
  }: Props
) =>
  <a
    style={style}
    href={`//${window.location.hostname}`}
  >
    <img
      alt='opctl logo'
      style={{
        height: '2rem'
      }}
      src="/logo.svg"
    />
  </a>
