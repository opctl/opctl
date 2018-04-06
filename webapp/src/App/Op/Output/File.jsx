import React from 'react';
import Description from '../Param/Description';

export default ({
                  name,
                  param,
                  opRef,
                  value,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <Description value={param.description} opRef={opRef}/>
      {
        value
        ?<a className='form-control' href={value || param.default || ''}>view file</a>
        : <input
        className='form-control'
        id={name}
        readOnly={true}
        type='text'
        value={value || param.default || ""}
      /> 
      }
    </div>
  );
}
