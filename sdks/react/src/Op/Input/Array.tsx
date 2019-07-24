import React from 'react'
import jsYaml from 'js-yaml'
import TextArea from './AceEditor'
import ParamArray from '@opctl/sdk/src/types/param/array'
import paramArrayValidate from '@opctl/sdk/src/opspec/interpreter/opcall/params/param/array/validate'

interface Props {
  array: ParamArray
  name: string
  onInvalid?: () => any | null | undefined
  onValid: (value: any) => any
  opRef: string
  value: any
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
    onValid={(value: any) => onValid(jsYaml.safeLoad(value))}
    opRef={opRef}
    validate={(value: any) => {
      try {
        return paramArrayValidate(jsYaml.safeLoad(value), array.constraints)
      } catch (err) {
        return [err]
      }
    }}
    value={value || jsYaml.safeDump(array.default ? array.default : '')}
  />
}
