import React from 'react'
import jsYaml from 'js-yaml'
import TextArea from './AceEditor'
import paramObjectValidate from '@opctl/sdk/lib/opspec/interpreter/opcall/params/param/object/validate'

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
          return paramObjectValidate(jsYaml.safeLoad(value), object.constraints)
        } catch (err) {
          return [err]
        }
      }}
      value={jsYaml.safeDump(value || object.default || '')}
    />
  )
