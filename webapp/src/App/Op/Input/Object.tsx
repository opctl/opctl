import React from 'react'
import jsYaml from 'js-yaml'
import TextArea from './AceEditor'
import opspecDataValidator from '@opspec/sdk/lib/data/object/validator'

interface Props {
  name: string
  object
  onInvalid
  onValid
  opRef: string
  value
}

export default (
  {
    name,
    object,
    onInvalid,
    onValid,
    opRef,
    value
  }: Props
) => (
    <TextArea
      description={object.description}
      name={name}
      onInvalid={onInvalid}
      onValid={value => onValid(jsYaml.safeLoad(value))}
      opRef={opRef}
      validate={value => {
        try {
          return opspecDataValidator.validate(jsYaml.safeLoad(value), object.constraints)
        } catch (err) {
          return [err]
        }
      }}
      value={value || jsYaml.safeDump(object.default ? object.default : '')}
    />
  )
