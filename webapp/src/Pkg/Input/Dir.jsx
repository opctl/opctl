import React from 'react';
import Input from './Input';

export default ({name, onValid, dir}) => (<Input
  description={dir.description}
  name={name}
  type='text'
  value={dir.default}
  // @TODO validate
  validate={value => ([])}
  onValid={value => onValid({dir: value})}
/>);

