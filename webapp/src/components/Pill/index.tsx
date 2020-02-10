import React from 'react'
import brandColors from '../../brandColors'
import { css, cx } from 'emotion'

interface Props {
  className?: string
  children: any
  onClick?: (e: any) => void
}

export default (
  {
    className,
    children,
    onClick
  }: Props
) =>
  <div
    className={cx(
      css({
        backgroundColor: brandColors.white,
        borderRadius: '1000vh',
        border: `thin solid`,
        color: brandColors.active,
        cursor: 'pointer',
        display: 'inline-block',
        padding: '.8rem 1.4rem',
        textAlign: 'center'
      }),
      className
    )}
    onClick={onClick} 
  >
    {
      children
    }
  </div>
