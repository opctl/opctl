import React from 'react';
import PkgRef from './PkgRef'
import Inputs from './Inputs'
import Outputs from './Outputs'

export default function Pkg(props) {
  const pkg = props.value;
  return (
    <div>
      <h1><PkgRef name={pkg.name} version={pkg.version}/></h1>
      <p className="lead">{pkg.description}</p>
      <Inputs value={pkg.inputs}/>
      <Outputs value={pkg.outputs}/>
    </div>
  );
}
