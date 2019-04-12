import React from 'react'
import jsYaml from 'js-yaml'
import TextArea from './AceEditor'
import paramArrayValidate from '@opctl/sdk/lib/opspec/interpreter/opcall/params/param/array/validate'

interface Props {
  array
  name: string
  onInvalid
  onValid
  opRef: string
  value
}

export default (
  {
    array,
    name,
    onInvalid,
    onValid,
    opRef,
    value
  }: Props
) => {
  return <TextArea
    description={array.description}
    name={name}
    onInvalid={onInvalid}
    onValid={value => onValid(jsYaml.safeLoad(value))}
    opRef={opRef}
    validate={value => {
      try {
        return paramArrayValidate(jsYaml.safeLoad(value), array.constraints)
      } catch (err) {
        return [err]
      }
    }}
    value={value || jsYaml.safeDump(array.default ? array.default : '')}
  />
}
