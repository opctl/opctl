import React from 'react'
import { css } from '@emotion/css'
import brandColors from '../../brandColors'

/**
 * A panel
 */
export default function Panel(props) {
  return (
    <div
      className={css({
        width: '100%',
        borderTop: `solid .1rem ${brandColors.reallyLightGray} !important`
      })}
    >
      {props.children}
    </div>
  )
}
