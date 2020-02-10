import React, { ReactElement } from 'react'
import { css, cx } from 'emotion'
import brandColors from '../../brandColors'

interface Props {
  className?: string
  children: ReactElement | ReactElement[]
}

/**
 * A row of the panel
 */
export default (
  {
    className,
    children
  }: Props
) =>
  <div
    className={cx(
      css({
        backgroundColor: brandColors.white,
        borderBottom: `solid thin ${brandColors.reallyLightGray} !important`,
        color: brandColors.black,
        padding: '1rem 1.2rem',
        width: '100%'
      }),
      className
    )}
  >
    {children}
  </div>
