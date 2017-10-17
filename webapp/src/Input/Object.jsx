import React from 'react';
import jsYaml from 'js-yaml';
import TextArea from './TextArea';
import opspecDataValidator from '@opspec/sdk/lib/data/object/validator';

export default ({name, onValid, object}) => (
  <TextArea
    description={object.description}
    name={name}
    value={jsYaml.safeDump(object.default ? object.default : '')}
    validate={value => opspecDataValidator.validate(jsYaml.safeLoad(value), object.constraints)}
    onChange={value => onValid({object: jsYaml.safeLoad(value)})}
  />
);
