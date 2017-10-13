import React from 'react';
import InputArray from './Array';
import InputDir from './Dir';
import InputFile from './File';
import InputNumber from './Number';
import InputObject from './Object';
import InputSocket from './Socket';
import InputString from './String';

export default (props) => {
  // delegate to component for input
  if (props.input.array) {
    return <InputArray
      array={props.input.array}
      name={props.name}
      onInvalid={props.onInvalid}
      onValid={props.onValid}
    />
  } else if (props.input.dir) {
    return <InputDir
      dir={props.input.dir}
      name={props.name}
      onInvalid={props.onInvalid}
      onValid={props.onValid}
    />
  } else if (props.input.file) {
    return <InputFile
      file={props.input.file}
      name={props.name}
      onInvalid={props.onInvalid}
      onValid={props.onValid}
    />
  } else if (props.input.number) {
    return <InputNumber
      onInvalid={props.onInvalid}
      onValid={props.onValid}
      name={props.name}
      number={props.input.number}
    />
  } else if (props.input.object) {
    return <InputObject
      name={props.name}
      object={props.input.object}
      onInvalid={props.onInvalid}
      onValid={props.onValid}
    />
  } else if (props.input.socket) {
    return <InputSocket
      name={props.name}
      onInvalid={props.onInvalid}
      onValid={props.onValid}
      socket={props.input.socket}
    />
  } else if (props.input.string) {
    return <InputString
      name={props.name}
      onInvalid={props.onInvalid}
      onValid={props.onValid}
      string={props.input.string}
    />
  }
  return null
}
