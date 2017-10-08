import React from 'react';

export default ({
                  string,
                  name,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{string.description}</p>
      <input
        className='form-control'
        defaultValue={string.default}
        id={name}
        readOnly={true}
        type='text'
      />
    </div>
  );
}
