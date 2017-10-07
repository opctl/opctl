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
  if (props.output.array) {
    return (<OutputArray name={props.name} array={props.output.array}/>);
  } else if (props.output.dir) {
    return (<OutputDir name={props.name} dir={props.output.dir}/>);
  } else if (props.output.file) {
    return (<OutputFile name={props.name} file={props.output.file}/>);
  } else if (props.output.number) {
    return (<OutputNumber name={props.name} number={props.output.number}/>);
  } else if (props.output.object) {
    return (<OutputObject name={props.name} object={props.output.object}/>);
  } else if (props.output.socket) {
    return (<OutputSocket name={props.name} socket={props.output.socket}/>);
  } else if (props.output.string) {
    return (<OutputString name={props.name} string={props.output.string}/>);
  }
  return null
}
