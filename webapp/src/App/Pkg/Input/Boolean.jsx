import React from 'react';
import Input from '../../Input';

export default ({name, onValid, pkgRef, boolean, value}) =>
  <Input
    description={boolean.description}
    name={name}
    onValid={value => onValid(value)}
    pkgRef={pkgRef}
    type='checkbox'
    validate={value => ([])}
    value={'undefined' === typeof value? boolean.default : value}
  />;