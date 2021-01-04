import React from 'react'
import { css } from '@emotion/css'
import brandColors from '../../brandColors'

export default function RightNav({
  isCompleteDisabled,
  onClick
}) {
  return (
    <div
      className={css({
        flex: '1 0 0',
        textAlign: 'right',
        marginRight: '1.2rem',
        ...isCompleteDisabled
          ? {
            color: brandColors.lightGray,
            cursor: 'default'
          }
          : {
            color: brandColors.active,
            cursor: 'pointer',
          }
      })}
    >
      {
        onClick &&
        <div
          {...!isCompleteDisabled && { onClick }}
        >
          Save
      </div>
      }
    </div>

  )
}