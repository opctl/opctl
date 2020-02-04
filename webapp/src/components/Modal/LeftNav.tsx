import React from 'react'
import { css } from 'emotion'
import brandColors from '../../brandColors'

export default ({ onClick }) => 
  <div
    className={css({
      color: brandColors.blue,
      cursor: 'pointer',
      flex: '1 0 0',
      marginLeft: '1.2rem'
    })}
  >
    {
      onClick &&
      <div
        onClick={onClick}
      >
        Close
      </div>
    }
  </div>
