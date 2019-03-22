import React from 'react'
import jsYaml from 'js-yaml'
import TextArea from './AceEditor'
import opspecDataValidator from '@opspec/sdk/lib/data/array/validator'

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
        return opspecDataValidator.validate(jsYaml.safeLoad(value), array.constraints)
      } catch (err) {
        return [err]
      }
    }}
    value={value || jsYaml.safeDump(array.default ? array.default : '')}
  />
}
