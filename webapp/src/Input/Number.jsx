import React from 'react';
import Input from './Input';
import opspecDataValidator from '@opspec/sdk/lib/data/number/validator';

export default ({name, onArgChange, number}) => (
  <Input
    description={number.description}
    name={name}
    type={number.isSecret ? 'password' : 'text'}
    value={number.default}
    validate={value => opspecDataValidator.validate(value, Object.assign({type: 'number'}, number.constraints))}
    onChange={value => onArgChange({number: value})}
  />
);
