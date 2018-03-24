import React from 'react';
import InputArray from './Array';
import InputBoolean from './Boolean';
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
                  opRef,
                  value,
                }) => {
  const environment = contentStore.get({key: 'environment'}) || [];
  const variable = environment.find(variable => variable.name === name) || {};
  value = value || variable.value;

  // delegate to component for input
  if (input.array) {
    return <InputArray
      array={input.array}
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      opRef={opRef}
      value={value}
    />
  } else if (input.boolean) {
    return <InputBoolean
      boolean={input.boolean}
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      opRef={opRef}
      value={value}
    />
  } else if (input.dir) {
    return <InputDir
      dir={input.dir}
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      opRef={opRef}
      value={value}
    />
  } else if (input.file) {
    return <InputFile
      file={input.file}
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      opRef={opRef}
      value={value}
    />
  } else if (input.number) {
    return <InputNumber
      onInvalid={onInvalid}
      onValid={onValid}
      name={name}
      number={input.number}
      opRef={opRef}
      value={value}
    />
  } else if (input.object) {
    return <InputObject
      name={name}
      object={input.object}
      onInvalid={onInvalid}
      onValid={onValid}
      opRef={opRef}
      value={value}
    />
  } else if (input.socket) {
    return <InputSocket
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      socket={input.socket}
      opRef={opRef}
      value={value}
    />
  } else if (input.string) {
    return <InputString
      name={name}
      onInvalid={onInvalid}
      onValid={onValid}
      string={input.string}
      opRef={opRef}
      value={value}
    />
  }
  return null
}
