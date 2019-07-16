import React from 'react'
import ModelParamBoolean from '@opctl/sdk/lib/model/param/boolean'
import _DomInput from './_DomInput'

interface Props {
  name: string
  onValid: (value: any) => any
  opRef: string
  boolean: ModelParamBoolean
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
