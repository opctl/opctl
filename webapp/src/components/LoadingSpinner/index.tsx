import React from 'react'
import brandColors from '../../brandColors'
import { ReactComponent as ProgressIcon } from '../../icons/Progress.svg'

export default function LoadingSpinner() {
  return (
    <div
      style={{
        display: 'flex',
        justifyContent: 'center',
        width: '100%'
      }}
    >
      <ProgressIcon
        style={{
          fill: brandColors.active,
          width: '2rem',
          height: '2rem'
        }}
      />
    </div>
  )
}
