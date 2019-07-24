import React from 'react'
import ParamFile from '@opctl/sdk/src/types/param/file'
import _DomInput from './_DomInput'

interface Props {
  file: ParamFile
  name: string
  onValid: (value: any) => any
  opRef: string
  value: any
}

export default (
  {
    file,
    name,
    onValid,
    opRef,
    value
  }: Props
) =>
  <_DomInput
    description={file.description}
    name={name}
    onValid={value => onValid(value)}
    opRef={opRef}
    type='text'
    // @TODO validate
    validate={() => ([])}
    value={value || file.default}
  />
