import React from 'react';
import OutputArray from './Array';
import OutputDir from './Dir';
import OutputFile from './File';
import OutputNumber from './Number';
import OutputObject from './Object';
import OutputSocket from './Socket';
import OutputString from './String';

export default (props) => {
  // delegate to component for output
  if (props.param.array) {
    return <OutputArray
      name={props.name}
      param={props.param.array}
      value={props.value.array || props.value.string || props.value.file}
    />
  } else if (props.param.dir) {
    return <OutputDir
      name={props.name}
      param={props.param.dir}
      value={props.value.dir}
    />
  } else if (props.param.file) {
    return <OutputFile
      name={props.name}
      param={props.param.file}
      value={props.value.file || props.value.string || props.value.number || props.value.array || props.value.object}
    />
  } else if (props.param.number) {
    return <OutputNumber
      name={props.name}
      param={props.param.number}
      value={props.value.number || props.value.file}
    />
  } else if (props.param.object) {
    return <OutputObject
      name={props.name}
      param={props.param.object}
      value={props.value.object || props.value.string || props.value.file}
    />
  } else if (props.param.socket) {
    return <OutputSocket
      name={props.name}
      param={props.param.socket}
      value={props.value.socket}
    />
  } else if (props.param.string) {
    return <OutputString
      name={props.name}
      param={props.param.string}
      value={props.value.string || props.value.number || props.value.array || props.value.object || props.value.file}
    />
  }
  return null
}
