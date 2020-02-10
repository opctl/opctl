import React, { CSSProperties } from 'react'
import { css, cx } from 'emotion'

interface Props {
  className?: string
  children
  style?: CSSProperties
}

/**
 * Blocks & Pads content
 */
export default ( 
  {
    className,
    style,
    children
  }: Props
) =>
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