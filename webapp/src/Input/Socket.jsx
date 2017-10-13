import React from 'react';
import Input from './Input';

export default ({name, onValid, socket}) => (<Input
  description={socket.description}
  name={name}
  type='text'
  // @TODO validate
  validate={value => ([])}
  onValid={value => onValid({socket: value})}
/>);

