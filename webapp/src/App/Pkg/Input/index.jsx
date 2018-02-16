import React from 'react';
import InputArray from './Array';
import InputDir from './Dir';
import InputFile from './File';
import InputNumber from './Number';
import InputObject from './Object';
import InputSocket from './Socket';
import InputString from './String';
import contentStore from '../../../core/contentStore';

export default ({
                  input,
                  name,
                  onInvalid,
                  onValid,
                  pkgRef,
                  value,
                }) => {
  const variable = contentStore.get({key: 'environment'}).find(variable => variable.name === name) || {};
  value = value || variable.value;

  // delegate to component for input
  if (input.array) {
    return <InputArray
      array={input.array}
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      pkgRef={pkgRef}
      value={value}
    />
  } else if (input.dir) {
    return <InputDir
      dir={input.dir}
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      pkgRef={pkgRef}
      value={value}
    />
  } else if (input.file) {
    return <InputFile
      file={input.file}
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      pkgRef={pkgRef}
      value={value}
    />
  } else if (input.number) {
    return <InputNumber
      onInvalid={onInvalid}
      onValid={onValid}
      name={name}
      number={input.number}
      pkgRef={pkgRef}
      value={value}
    />
  } else if (input.object) {
    return <InputObject
      name={name}
      object={input.object}
      onInvalid={onInvalid}
      onValid={onValid}
      pkgRef={pkgRef}
      value={value}
    />
  } else if (input.socket) {
    return <InputSocket
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      socket={input.socket}
      pkgRef={pkgRef}
      value={value}
    />
  } else if (input.string) {
    return <InputString
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      string={input.string}
      pkgRef={pkgRef}
      value={value}
    />
  }
  return null
}
