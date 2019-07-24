import React from 'react'
import _DomInput from './_DomInput'
import Select from './Select'
import ParamString from '@opctl/sdk/src/types/param/string'
import paramStringValidate from '@opctl/sdk/src/opspec/interpreter/opcall/params/param/string/validate'

interface Props {
  name: string
  onInvalid?: () => any | null | undefined
  onValid: (value: any) => any
  opRef: string
  string: ParamString
  value: any
}

export default (
  {
    name,
    onInvalid,
    onValid,
    opRef,
    string,
    value
  }: Props
) => {
  if (string.constraints && !string.isSecret && string.constraints.enum) {
    return <Select
      description={string.description}
      name={name}
      options={string.constraints.enum.map((item:any) => 
        ({ name: item, value: item })
      )}
      onInvalid={onInvalid}
      onValid={value => onValid(value)}
      opRef={opRef}
      validate={value => paramStringValidate(value, string.constraints)}
      value={value || string.default}
    />
  }
  return <_DomInput
    description={string.description}
    name={name}
    onInvalid={onInvalid}
    onValid={value => onValid(value)}
    opRef={opRef}
    type={string.isSecret ? 'password' : 'text'}
    validate={value => paramStringValidate(value, string.constraints)}
    value={value || string.default}
  />
}
