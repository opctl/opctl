import React from 'react'
import Param from '@opctl/sdk/src/types/param'
import OutputArray from './Array'
import OutputBoolean from './Boolean'
import OutputDir from './Dir'
import OutputFile from './File'
import OutputNumber from './Number'
import OutputObject from './Object'
import OutputSocket from './Socket'
import OutputString from './String'

interface Props {
  apiBaseUrl: string
  name: string
  param: Param
  opRef: string
  value: any
}

export default (
  {
    apiBaseUrl,
    name,
    param,
    opRef,
    value
  }: Props
) => {
  // delegate to component for output
  if (param.array) {
    return <OutputArray
      apiBaseUrl={apiBaseUrl}
      name={name}
      param={param.array}
      opRef={opRef}
      value={value.array || value.string || value.file}
    />
  } else if (param.boolean) {
    return <OutputBoolean
      apiBaseUrl={apiBaseUrl}
      name={name}
      param={param.boolean}
      opRef={opRef}
      value={value.boolean}
    />
  } else if (param.dir) {
    return <OutputDir
      name={name}
      param={param.dir}
      opRef={opRef}
      value={value.dir}
    />
  } else if (param.file) {
    return <OutputFile
      name={name}
      param={param.file}
      opRef={opRef}
      value={value.file || value.string || value.number || value.array || value.object}
    />
  } else if (param.number) {
    return <OutputNumber
      apiBaseUrl={apiBaseUrl}
      name={name}
      param={param.number}
      opRef={opRef}
      value={value.number || value.file}
    />
  } else if (param.object) {
    return <OutputObject
      apiBaseUrl={apiBaseUrl}
      name={name}
      param={param.object}
      opRef={opRef}
      value={value.object || value.string || value.file}
    />
  } else if (param.socket) {
    return <OutputSocket
      name={name}
      param={param.socket}
      opRef={opRef}
    />
  } else if (param.string) {
    return <OutputString
      apiBaseUrl={apiBaseUrl}
      name={name}
      param={param.string}
      opRef={opRef}
      value={value.string || value.number || value.array || value.object || value.file}
    />
  }
  return null
}
