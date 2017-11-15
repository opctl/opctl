import React from 'react';
import Input from './Input';

export default ({dir, name, onValid, pkgRef}) =>
  <Input
    description={dir.description}
    name={name}
    onValid={value => onValid({dir: value})}
    pkgRef={pkgRef}
    type='text'
    value={dir.default}
    // @TODO validate
    validate={value => ([])}
  />;

