import React from 'react';
import Input from './Input/index';

export default ({value, onInvalid, onValid}) => {
  const inputs = Object.entries(value || {}).map(([name, input]) => (
    <Input
      onInvalid={() => (onInvalid(name))}
      onValid={value => (onValid(name, value))}
      name={name}
      input={input}
      key={name}
    />
  ));

  return (
    <div>
      <h2>Inputs:</h2>
      {inputs}
    </div>
  );
}
