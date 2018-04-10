import React from 'react'
import Input from '../../Input'

export default ({name, onValid, opRef, boolean, value}) =>
  <Input
    description={boolean.description}
    name={name}
    onValid={value => onValid(value)}
    opRef={opRef}
    type='checkbox'
    validate={value => ([])}
    value={typeof value === 'undefined' ? boolean.default : value}
  />
