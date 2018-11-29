import React from 'react'
import Input from '../../Input'

export default ({ dir, name, onValid, opRef, value }) =>
  <Input
    description={dir.description}
    name={name}
    onValid={value => onValid(value)}
    opRef={opRef}
    type='text'
    // @TODO validate
    validate={value => ([])}
    value={value || dir.default}
  />
