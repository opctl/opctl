import React from 'react'
import ModelParamDir from '@opctl/sdk/lib/model/param/dir'
import _DomInput from './_DomInput'

interface Props {
  dir: ModelParamDir
  name: string
  onValid: (value: any) => any
  opRef: string
  value: any
}

export default (
  {
    dir,
    name,
    onValid,
    opRef,
    value
  }: Props
) =>
  <_DomInput
    description={dir.description}
    name={name}
    onValid={value => onValid(value)}
    opRef={opRef}
    type='text'
    // @TODO validate
    validate={value => ([])}
    value={value || dir.default}
  />
