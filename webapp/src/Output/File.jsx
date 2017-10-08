import React from 'react';

export default ({
                  file,
                  name,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{file.description}</p>
      <input
        className='form-control'
        defaultValue={file.default}
        id={name}
        readOnly={true}
        type='text'
      />
    </div>
  );
}
