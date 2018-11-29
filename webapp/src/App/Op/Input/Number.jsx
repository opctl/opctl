import React from 'react'
import Input from '../../Input'
import Select from './Select'
import opspecDataValidator from '@opspec/sdk/lib/data/number/validator'

export default ({ name, number, onInvalid, onValid, opRef, value }) => {
  if (number.constraints && !number.isSecret && number.constraints.enum) {
    return <Select
      description={number.description}
      name={name}
      onInvalid={onInvalid}
      onValid={value => onValid(value)}
      options={number.constraints.enum.map(item => ({ name: item, value: item }))}
      opRef={opRef}
      validate={value => opspecDataValidator.validate(Number(value), number.constraints)}
      value={value || number.default}
    />
  }
  return <Input
    description={number.description}
    name={name}
    onInvalid={onInvalid}
    onValid={value => onValid(value)}
    opRef={opRef}
    type={number.isSecret ? 'password' : 'number'}
    validate={value => opspecDataValidator.validate(Number(value), number.constraints)}
    value={value || number.default}
  />
}
