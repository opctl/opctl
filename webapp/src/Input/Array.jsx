import React from 'react';
import jsYaml from 'js-yaml';
import TextArea from './TextArea';
import opspecDataValidator from '@opspec/sdk/lib/data/array/validator';

export default ({name, onValid, array}) => {
  return <TextArea
    description={array.description}
    name={name}
    value={jsYaml.safeDump(array.default ? array.default : '')}
    validate={value => opspecDataValidator.validate(jsYaml.safeLoad(value), array.constraints)}
    onChange={value => onValid({array: jsYaml.safeLoad(value)})}
  />
}
