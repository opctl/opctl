import React from 'react';
import PkgRef from './PkgRef'
import Inputs from './Inputs'
import Outputs from './Outputs'

export default function Pkg(props) {
  return (
    <div>
      <h1><PkgRef name={props.name} version={props.version}/></h1>
      <p className="lead">{props.description}</p>
      <Inputs value={props.inputs}/>
      <Outputs value={props.outputs}/>
    </div>
  );
}
