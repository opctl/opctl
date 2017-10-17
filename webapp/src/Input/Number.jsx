import React from 'react';
import Input from './Input';
import Select from './Select';
import opspecDataValidator from '@opspec/sdk/lib/data/number/validator';

export default ({name, onInvalid, onValid, number}) => {
  if (number.constraints && !number.isSecret && number.constraints.enum) {
    return <Select
      description={number.description}
      name={name}
      onInvalid={onInvalid}
      onValid={value => onValid({number: value})}
      options={number.constraints.enum.map(item => ({name: item, value: item}))}
      validate={value => opspecDataValidator.validate(Number(value), number.constraints)}
      value={number.default}
    />
  }
  return <Input
    description={number.description}
    name={name}
    onInvalid={onInvalid}
    onValid={value => onValid({number: value})}
    type={number.isSecret ? 'password' : 'number'}
    validate={value => opspecDataValidator.validate(Number(value), number.constraints)}
    value={number.default}
  />
};
