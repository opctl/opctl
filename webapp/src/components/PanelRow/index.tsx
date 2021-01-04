import React, { ReactElement } from 'react'
import { css, cx } from '@emotion/css'
import brandColors from '../../brandColors'

interface Props {
  className?: string
  children: ReactElement | ReactElement[]
}

/**
 * A row of the panel
 */
export default function PanelRow(
  {
    className,
    children
  }: Props
) {
  return (
    <div
      className={cx(
        css({
          backgroundColor: brandColors.white,
          borderBottom: `solid .1rem ${brandColors.reallyLightGray} !important`,
          color: brandColors.black,
          padding: '1rem 1.2rem',
          width: '100%'
        }),
        className
      )}
    >
      {children}
    </div>

  )
}