import React from 'react';
import OutputArray from './Array';
import OutputDir from './Dir';
import OutputFile from './File';
import OutputNumber from './Number';
import OutputObject from './Object';
import OutputSocket from './Socket';
import OutputString from './String';

export default ({
                  name,
                  param,
                  pkgRef,
                  value,
                }) => {
  // delegate to component for output
  if (param.array) {
    return <OutputArray
      name={name}
      param={param.array}
      pkgRef={pkgRef}
      value={value.array || value.string || value.file}
    />
  } else if (param.dir) {
    return <OutputDir
      name={name}
      param={param.dir}
      pkgRef={pkgRef}
      value={value.dir}
    />
  } else if (param.file) {
    return <OutputFile
      name={name}
      param={param.file}
      pkgRef={pkgRef}
      value={value.file || value.string || value.number || value.array || value.object}
    />
  } else if (param.number) {
    return <OutputNumber
      name={name}
      param={param.number}
      pkgRef={pkgRef}
      value={value.number || value.file}
    />
  } else if (param.object) {
    return <OutputObject
      name={name}
      param={param.object}
      pkgRef={pkgRef}
      value={value.object || value.string || value.file}
    />
  } else if (param.socket) {
    return <OutputSocket
      name={name}
      param={param.socket}
      pkgRef={pkgRef}
      value={value.socket}
    />
  } else if (param.string) {
    return <OutputString
      name={name}
      param={param.string}
      pkgRef={pkgRef}
      value={value.string || value.number || value.array || value.object || value.file}
    />
  }
  return null
}
