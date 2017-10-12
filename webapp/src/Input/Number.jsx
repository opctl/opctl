import React from 'react';
import Input from './Input';
import Select from './Select';
import opspecDataValidator from '@opspec/sdk/lib/data/number/validator';

export default ({name, onArgChange, number}) => {
  if (number.constraints && !number.isSecret && number.constraints.enum) {
    return <Select
      description={number.description}
      name={name}
      value={number.default}
      options={number.constraints.enum.map(item => ({name: item, value: item}))}
      onChange={value => onArgChange({number: value})}
    />
  }
  return <Input
    description={number.description}
    name={name}
    type={number.isSecret ? 'password' : 'text'}
    value={number.default}
    validate={value => opspecDataValidator.validate(value, Object.assign({type: 'number'}, number.constraints))}
    onChange={value => onArgChange({number: value})}
  />
};
