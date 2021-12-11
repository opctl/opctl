import React from 'react'
import jsYaml from 'js-yaml'
import TextArea from './AceEditor'
import {ErrorObject} from 'ajv'
import ModelParamObject from '@opctl/sdk/src/model/param/object'
import paramObjectValidate from '@opctl/sdk/src/opspec/interpreter/opcall/params/param/object/validate'

interface Props {
  name: string
  object: ModelParamObject
  onInvalid?: () => any | null | undefined
  onValid: (value: any) => any
  opRef: string
  value: any
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
          return [err as ErrorObject]
        }
      }}
      value={value
        ? jsYaml.safeDump(value)
        : jsYaml.safeDump(object.default || '')
      }
    />
  )
