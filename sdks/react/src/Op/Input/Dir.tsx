import React from 'react'
import ParamDir from '@opctl/sdk/src/types/param/dir'
import _DomInput from './_DomInput'

interface Props {
  dir: ParamDir
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
