import React from 'react';
import jsYaml from 'js-yaml';
import Textarea from 'react-textarea-autosize';

export default ({
                  object,
                  name,
                }) => {
  return (
    <div className='form-group'>
      <label className='form-control-label' htmlFor={name}>{name}</label>
      <p className='custom-control-description'>{object.description}</p>
      <Textarea
        className='form-control'
        defaultValue={jsYaml.safeDump(object.default)}
        id={name}
        readOnly={true}
      />
    </div>
  );
}
