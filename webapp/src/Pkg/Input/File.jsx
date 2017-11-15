import React from 'react';
import Input from './Input';

export default ({file, name, onValid, pkgRef}) =>
  <Input
    description={file.description}
    name={name}
    onValid={value => onValid({file: value})}
    pkgRef={pkgRef}
    type='text'
    value={file.default}
    // @TODO validate
    validate={() => ([])}
  />;

