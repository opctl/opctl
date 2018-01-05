import React from 'react';
import Input from './Input';

export default ({name, onValid, pkgRef, socket, value}) =>
  <Input
    description={socket.description}
    name={name}
    onValid={value => onValid(value)}
    pkgRef={pkgRef}
    type='text'
    // @TODO validate
    validate={value => ([])}
    value={value}
  />;

