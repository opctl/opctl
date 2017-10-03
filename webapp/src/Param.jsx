import React from 'react';
import ParamDir from './Param_Dir';
import ParamFile from './Param_File';
import ParamNumber from './Param_Number';
import ParamSocket from './Param_Socket';
import ParamString from './Param_String';

export default function Param(props){
    // delegate to component for param
    if (props.param.dir){
        return (<ParamDir name={props.name} dir={props.param.dir} />);
    } else if (props.param.file){
        return (<ParamFile name={props.name} file={props.param.file} />);
    } else if (props.param.number){
        return (<ParamNumber name={props.name} number={props.param.number} />);
    } else if (props.param.socket){
        return (<ParamSocket name={props.name} socket={props.param.socket} />);
    } else if (props.param.string) {
        return (<ParamString name={props.name} string={props.param.string} />);
    }
}
