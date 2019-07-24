import React from 'react'
import ParamSocket from '@opctl/sdk/src/types/param/socket'
import _DomInput from './_DomInput'

interface Props {
  name: string
  onValid: (value: any) => any
  opRef: string
  socket: ParamSocket
  value: any
}

export default (
  {
    name,
    onValid,
    opRef,
    socket,
    value
  }: Props
) =>
  <_DomInput
    description={socket.description}
    name={name}
    onValid={value => onValid(value)}
    opRef={opRef}
    type='text'
    // @TODO validate
    validate={value => ([])}
    value={value}
  />
