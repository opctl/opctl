import React from 'react';
import jsYaml from 'js-yaml';
import TextArea from './TextArea';
import opspecDataValidator from '@opspec/sdk/lib/data/object/validator';

export default ({name, object, onInvalid, onValid}) => (
  <TextArea
    description={object.description}
    name={name}
    onInvalid={onInvalid}
    onValid={value => onValid({object: jsYaml.safeLoad(value)})}
    value={jsYaml.safeDump(object.default ? object.default : '')}
    validate={value => opspecDataValidator.validate(jsYaml.safeLoad(value), object.constraints)}
  />
);
