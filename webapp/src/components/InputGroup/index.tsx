import React from 'react'
import { css, cx } from '@emotion/css'

/**
 * A row of inputs
 */
export default function InputGroup(props) {
  return (
    <div
      className={cx(
        css({
          width: '100%',
          marginBottom: '1rem'
        }),
        props.className
      )}
    >
      {props.children}
    </div>
  )
}
