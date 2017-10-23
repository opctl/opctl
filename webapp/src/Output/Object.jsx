import React from 'react';
import jsYaml from 'js-yaml';
import Textarea from 'react-textarea-autosize';

export default ({
                  name,
                  param,
                  value,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{param.description}</p>
      <Textarea
        className='form-control'
        defaultValue={value || jsYaml.safeDump(param.default)}
        id={name}
        readOnly={true}
      />
    </div>
  );
}
