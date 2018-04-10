import React from 'react'
import Input from '../../Input'

export default ({name, onValid, opRef, socket, value}) =>
  <Input
    description={socket.description}
    name={name}
    onValid={value => onValid(value)}
    opRef={opRef}
    type='text'
    // @TODO validate
    validate={value => ([])}
    value={value}
  />
