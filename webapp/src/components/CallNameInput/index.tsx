import React from 'react'
import TextInput from '../TextInput'

interface Props {
    onChange: (value: string) => any
    value: string
}

export default (
    {
        onChange,
        value
    }: Props
) => <TextInput
        onChange={onChange}
        value={value}
        placeholder='Name'
    />  