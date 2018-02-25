import React from 'react';
import Input from '../../Input';

export default ({dir, name, onValid, pkgRef, value}) =>
  <Input
    description={dir.description}
    name={name}
    onValid={value => onValid(value)}
    pkgRef={pkgRef}
    type='text'
    // @TODO validate
    validate={value => ([])}
    value={value || dir.default}
  />;

