import React from 'react'
import Input from '../../Input'
import Select from './Select'
import paramStringValidate from '@opctl/sdk/lib/opspec/interpreter/opcall/params/param/string/validate'

export default ({ name, onInvalid, onValid, opRef, string, value }) => {
  if (string.constraints && !string.isSecret && string.constraints.enum) {
    return <Select
      description={string.description}
      name={name}
      options={string.constraints.enum.map(item => ({ name: item, value: item }))}
      onInvalid={onInvalid}
      onValid={value => onValid(value)}
      opRef={opRef}
      validate={value => paramStringValidate(value, string.constraints)}
      value={value || string.default}
    />
  }
  return <Input
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
