import React from 'react';
import jsYaml from 'js-yaml';
import TextArea from './AceEditor';
import opspecDataValidator from '@opspec/sdk/lib/data/object/validator';

export default ({name, object, onInvalid, onValid, pkgRef, value}) => (
  <TextArea
    description={object.description}
    name={name}
    onInvalid={onInvalid}
    onValid={value => onValid({object: jsYaml.safeLoad(value)})}
    pkgRef={pkgRef}
    validate={value => {
      try {
        return opspecDataValidator.validate(jsYaml.safeLoad(value), object.constraints)
      } catch (err) {
        return [err];
      }
    }}
    value={value || jsYaml.safeDump(object.default ? object.default : '')}
  />
);
