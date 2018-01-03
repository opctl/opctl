import React from 'react';
import Input from './Input';

export default ({file, name, onValid, pkgRef, value}) =>
  <Input
    description={file.description}
    name={name}
    onValid={value => onValid({file: value, value})}
    pkgRef={pkgRef}
    type='text'
    // @TODO validate
    validate={() => ([])}
    value={value || file.default}
  />;

