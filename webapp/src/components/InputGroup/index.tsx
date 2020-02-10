import React from 'react'
import { css, cx } from 'emotion'

/**
 * A row of inputs
 */
export default props =>
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
