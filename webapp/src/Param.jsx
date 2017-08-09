import React from 'react';
import ParamDir from './ParamDir';
import ParamFile from './ParamFile';
import ParamNumber from './ParamNumber';
import ParamSocket from './ParamSocket';
import ParamString from './ParamString';

export default function Param(props){
    // dereference known types
    const dirParam = props.type.dir;
    const fileParam = props.type.file;
    const numberParam = props.type.number;
    const socketParam = props.type.socket;
    const stringParam = props.type.string;

    // delegate to component for type
    if (dirParam){
        return (<ParamDir name={props.name} param={dirParam} />);
    } else if (fileParam){
        return (<ParamFile name={props.name} param={fileParam} />);
    } else if (numberParam){
        return (<ParamNumber name={props.name} param={numberParam} />);
    } else if (socketParam){
        return (<ParamSocket name={props.name} param={socketParam} />);
    } else if (stringParam) {
        return (<ParamString name={props.name} param={stringParam} />);
    }
}