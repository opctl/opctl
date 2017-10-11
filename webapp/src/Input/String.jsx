import React from 'react';
import Input from './Input';
import opspecDataValidator from '@opspec/sdk/lib/data/string/validator';

export default ({name, onArgChange, string}) => (
  <Input
    description={string.description}
    name={name}
    type={string.isSecret ? 'password' : 'text'}
    value={string.default}
    validate={value => opspecDataValidator.validate(value, string.constraints)}
    onChange={value => onArgChange({string: value})}
  />
);
