import React from 'react'
import ParamBoolean from '@opctl/sdk/src/types/param/boolean'
import _DomInput from './_DomInput'

interface Props {
  name: string
  onValid: (value: any) => any
  opRef: string
  boolean: ParamBoolean
  value: any
}

export default (
  {
    name,
    onValid,
    opRef,
    boolean,
    value
  }: Props
) =>
  <_DomInput
    description={boolean.description}
    name={name}
    onValid={(value: any) => onValid(value)}
    opRef={opRef}
    type='checkbox'
    validate={() => ([])}
    value={typeof value === 'undefined' ? boolean.default : value}
  />
