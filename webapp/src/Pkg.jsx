import React from 'react';
import PkgRef from './PkgRef'
import Inputs from './Inputs'
import Outputs from './Outputs'

export default ({value, pkgRef}) => {
  return (
    <div>
      <h1><PkgRef name={value.name} version={value.version}/></h1>
      <p className="lead">{value.description}</p>
      <Inputs value={value.inputs} pkgRef={pkgRef}/>
      <Outputs value={value.outputs}/>
    </div>
  );
}
