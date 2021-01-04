import React from 'react'
import brandColors from '../../brandColors'
import { css, cx } from '@emotion/css'

interface Props {
  className?: string
  children: any
  onClick?: (e: any) => void
}

export default function Pill(
  {
    className,
    children,
    onClick
  }: Props
) {
  return (
    <div
      className={cx(
        css({
          backgroundColor: brandColors.white,
          borderRadius: '1000vh',
          border: `.1rem solid`,
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

  )
}