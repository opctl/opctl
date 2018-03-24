import React from 'react';
import Input from '../../Input';

export default ({file, name, onValid, opRef, value}) =>
  <Input
    description={file.description}
    name={name}
    onValid={value => onValid(value)}
    opRef={opRef}
    type='text'
    // @TODO validate
    validate={() => ([])}
    value={value || file.default}
  />;

