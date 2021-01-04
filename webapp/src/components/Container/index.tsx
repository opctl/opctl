import React, { CSSProperties } from 'react'
import { css, cx } from '@emotion/css'

interface Props {
  className?: string
  children
  style?: CSSProperties
}

/**
 * Blocks & Pads content
 */
export default function Container(
  {
    className,
    style,
    children
  }: Props
) {
  return (
    <div
      className={
        cx(
          css({
            padding: '1.2rem'
          }),
          className
        )
      }
      style={style}
    >
      {
        children
      }
    </div>
  )
}