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
    return (<InputArray onChange={props.onChange} name={props.name} array={props.input.array}/>);
  } else if (props.input.dir) {
    return (<InputDir onChange={props.onChange} name={props.name} dir={props.input.dir}/>);
  } else if (props.input.file) {
    return (<InputFile onChange={props.onChange} name={props.name} file={props.input.file}/>);
  } else if (props.input.number) {
    return (<InputNumber onChange={props.onChange}  name={props.name} number={props.input.number}/>);
  } else if (props.input.object) {
    return (<InputObject onChange={props.onChange} name={props.name} object={props.input.object}/>);
  } else if (props.input.socket) {
    return (<InputSocket onChange={props.onChange} name={props.name} socket={props.input.socket}/>);
  } else if (props.input.string) {
    return (<InputString onChange={props.onChange} name={props.name} string={props.input.string}/>);
  }
  return null
}
