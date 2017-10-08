import React from 'react';
import jsYaml from 'js-yaml';
import Textarea from 'react-textarea-autosize';

export default ({
                  array,
                  name,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{array.description}</p>
      <Textarea
        className='form-control'
        defaultValue={jsYaml.safeDump(array.default)}
        id={name}
        readOnly={true}
      />
    </div>
  );
}
