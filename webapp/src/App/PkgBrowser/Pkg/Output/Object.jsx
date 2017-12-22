import React from 'react';
import jsYaml from 'js-yaml';
import Textarea from 'react-textarea-autosize';
import Description from './Description';

export default ({
                  name,
                  param,
                  pkgRef,
                  value,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <Description value={param.description} pkgRef={pkgRef}/>
      <Textarea
        className='form-control'
        defaultValue={value || jsYaml.safeDump(param.default)}
        id={name}
        readOnly={true}
      />
    </div>
  );
}
