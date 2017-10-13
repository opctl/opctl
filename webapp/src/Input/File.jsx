import React from 'react';
import Input from './Input';

export default ({name, onValid, file}) => (<Input
  description={file.description}
  name={name}
  type='text'
  value={file.default}
  // @TODO validate
  validate={value => ([])}
  onValid={value => onValid({file: value})}
/>);

