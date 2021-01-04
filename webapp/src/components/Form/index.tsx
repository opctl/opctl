import React from 'react'
import InputGroup from '../InputGroup'
import Pill from '../Pill'
import { css, cx } from '@emotion/css'

interface Props {
  className?: string
  children
  isValid?: boolean
  onSubmit?: () => void
  submitName?: string
}

export default function Form(
  {
    className,
    children,
    isValid,
    onSubmit,
    submitName
  }: Props
) {
  return (
    <form
      className={
        cx(
          css({
            padding: '1.2rem'
          }),
          className
        )}
      onSubmit={e => e.preventDefault()}
    >
      {children}
      {
        <input
          // include hidden submit to ensure enter always triggers submit
          // without submit input, browsers disable submit on enter
          type='submit'
          style={{ display: 'none' }}
        />
      }
      {
        onSubmit
          ? <InputGroup
            className={
              css({
                paddingTop: '1rem',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'flex-end'
              })
            }>
            <Pill
              className={css({
                ...!isValid
                  ? {
                    cursor: 'default',
                    opacity: .5
                  }
                  : {}
              })}
              onClick={
                isValid
                  ? onSubmit
                  : () => { }
              }
            >
              {submitName || 'Save'}
            </Pill>
          </InputGroup>
          : null
      }
    </form>

  )
}