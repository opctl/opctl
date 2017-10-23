import React from 'react';

export default ({
                  name,
                  param,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{param.description}</p>
      <input
        className='form-control'
        id={name}
        readOnly={true}
        type='text'
      />
    </div>
  );
}
