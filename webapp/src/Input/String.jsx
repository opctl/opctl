import React from 'react';
import Input from './Input';
import Select from './Select';
import opspecDataValidator from '@opspec/sdk/lib/data/string/validator';

export default ({name, onArgChange, string}) => {
  if (string.constraints && !string.isSecret && string.constraints.enum) {
    return <Select
      description={string.description}
      name={name}
      value={string.default}
      options={string.constraints.enum.map(item => ({name: item, value: item}))}
      onChange={value => onArgChange({string: value})}
    />
  }
  return <Input
    description={string.description}
    name={name}
    type={string.isSecret ? 'password' : 'text'}
    value={string.default}
    validate={value => opspecDataValidator.validate(value, Object.assign({type: 'string'}, string.constraints))}
    onChange={value => onArgChange({string: value})}
  />
};
