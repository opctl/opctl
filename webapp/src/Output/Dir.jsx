import React from 'react';

export default ({
                  dir,
                  name,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{dir.description}</p>
      <input
        className='form-control'
        defaultValue={dir.default}
        id={name}
        readOnly={true}
        type='text'
      />
    </div>
  );
}
