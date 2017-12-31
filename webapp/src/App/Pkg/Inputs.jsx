import React from 'react';
import Input from './Input/index';

export default ({value, onInvalid, onValid, pkgRef, values}) => {
  const inputs = Object.entries(value || {}).map(([name, input]) => (
    <Input
      onInvalid={() => (onInvalid(name))}
      onValid={value => (onValid(name, value))}
      name={name}
      pkgRef={pkgRef}
      input={input}
      key={name}
      value={values[name]|| null}
    />
  ));

  return (
    <div>
      <h2>Inputs:</h2>
      {inputs}
    </div>
  );
}
