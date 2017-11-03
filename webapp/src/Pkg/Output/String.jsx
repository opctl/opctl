import React from 'react';
import Description from './Description';

export default ({
                  name,
                  param,
                  value,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <Description value={param.description}/>
      <input
        className='form-control'
        id={name}
        readOnly={true}
        type={param.isSecret ? 'password' : 'text'}
        value={value || param.default || ""}
      />
    </div>
  );
}
