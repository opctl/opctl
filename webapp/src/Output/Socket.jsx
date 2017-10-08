import React from 'react';

export default ({
                  socket,
                  name,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{socket.description}</p>
      <input
        className='form-control'
        id={name}
        readOnly={true}
        type='text'
      />
    </div>
  );
}
