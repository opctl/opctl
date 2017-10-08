import React from 'react';

export default ({
                  number,
                  name,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{number.description}</p>
      <input
        className='form-control'
        defaultValue={number.default}
        id={name}
        readOnly={true}
        type='number'
      />
    </div>
  );
}
