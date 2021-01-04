import React from 'react'
import { cx, css } from '@emotion/css'
import brandColors from '../../brandColors'
import formElementStyles from '../formElementStyles'

interface Props {
  className?: string
  disabled?: boolean
  onBlur?
  onChange: (value: string) => any
  onClick?
  onKeyDown?
  placeholder?: string
  readOnly?: boolean
  value?: string
}

export default function TextInput(
  props: Props
) {
  return (
    <input
      {...props}
      className={cx(
        css(formElementStyles),
        css({
          backgroundColor: brandColors.white,
          color: brandColors.black,
          padding: '1rem 1.2rem',
          border: `solid .1rem ${brandColors.lightGray} !important`,
          width: '100%'
        }),
        props.className
      )}
      onChange={({ target }) => props.onChange(target.value)}
    />
  )
}
