import React from 'react';
import Input from './Input';

export default ({name, onValid, pkgRef, socket}) =>
  <Input
    description={socket.description}
    name={name}
    onValid={value => onValid({socket: value})}
    pkgRef={pkgRef}
    type='text'
    // @TODO validate
    validate={value => ([])}
  />;

