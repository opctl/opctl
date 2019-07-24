import React from 'react'
import _DomInput from './_DomInput'
import Select from './Select'
import ParamNumber from '@opctl/sdk/src/types/param/number'
import paramNumberValidate from '@opctl/sdk/src/opspec/interpreter/opcall/params/param/number/validate'

interface Props {
  name: string
  number: ParamNumber
  onInvalid?: () => any | null | undefined
  onValid: (value: any) => any
  opRef: string
  value: any
}

export default (
  {
    name,
    number,
    onInvalid,
    onValid,
    opRef,
    value
  }: Props
) => {
  if (number.constraints && !number.isSecret && number.constraints.enum) {
    return <Select
      description={number.description}
      name={name}
      onInvalid={onInvalid}
      onValid={value => onValid(value)}
      options={number.constraints.enum.map((item:any) =>
        ({ name: item, value: item }))
      }
      opRef={opRef}
      validate={value => paramNumberValidate(Number(value), number.constraints)}
      value={value || number.default}
    />
  }
  return <_DomInput
    description={number.description}
    name={name}
    onInvalid={onInvalid}
    onValid={value => onValid(value)}
    opRef={opRef}
    type={number.isSecret ? 'password' : 'number'}
    validate={value => paramNumberValidate(Number(value), number.constraints)}
    value={value || number.default}
  />
}
