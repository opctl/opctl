import React from 'react'
import TextInput from '../TextInput'

interface Props {
  onChange: (value: string) => any
  value: string
}

export default function CallNameInput(
  {
    onChange,
    value
  }: Props
) {
  return (
    <TextInput
      onChange={onChange}
      value={value}
      placeholder='Name'
    />
  )
} 