import React from 'react'
import { css } from 'emotion'
import brandColors from '../../brandColors'

export default ({
  isCompleteDisabled,
  onClick
}) =>
  <div
    className={css({
      color: brandColors.active,
      cursor: 'pointer',
      flex: '1 0 0',
      textAlign: 'right',
      marginRight: '1.2rem',
      ...isCompleteDisabled && {
        color: brandColors.lightGray,
        cursor: 'default'
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
